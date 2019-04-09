// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 5120
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub        *Hub
	sqlMessage *SqlMessage

	is_login bool

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	user_id int32
}
type LoginInfo struct {
	client  *Client
	user_id int32
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, origin_message, err := c.conn.ReadMessage()
		if err != nil {
			//log.Printf("my error: %v", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		myWSMessage := &WSMessage{}
		//解码数据
		proto.Unmarshal(origin_message, myWSMessage)

		//todo,buyer_auth 或 shop_auth 不应该由客户端传，应该放token里面。
		if myWSMessage.Type == "buyer_auth" {
			AuthMessage := &AuthMessage{}
			proto.Unmarshal(myWSMessage.Content, AuthMessage)
			//认证逻辑处理
			hmacSampleSecret := []byte("88888888")

			token, err := jwt.Parse(AuthMessage.Token, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return hmacSampleSecret, nil
			})

			if nil != err {
				log.Printf("auth err %v login ", err)

				//token 验证不通过，断开连接
				if _, ok := c.hub.clients[c]; ok {
					break
				}
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				//验证通过，加入 login
				temp_client_id, _ := strconv.ParseInt(claims["client_id"].(string), 10, 64)
				client_id := int32(temp_client_id)


				login_info := &LoginInfo{client: c, user_id: client_id}
				c.hub.login <- login_info
				c.is_login = true
				c.user_id = client_id
				log.Printf("buyer %v login success",client_id)
			} else {
				//token 验证不通过，断开连接
				if _, ok := c.hub.clients[c]; ok {
					break
				}
				fmt.Println(err)
			}
		} else if myWSMessage.Type == "shoper_auth" {
			AuthMessage := &AuthMessage{}
			proto.Unmarshal(myWSMessage.Content, AuthMessage)
			//店铺登录
			hmacSampleSecret := []byte("88888888")

			token, err := jwt.Parse(AuthMessage.Token, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return hmacSampleSecret, nil
			})
			if nil != err {
				log.Printf("auth err %v login ", err)
				//token 验证不通过，断开连接
				if _, ok := c.hub.clients[c]; ok {
					break
				}
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				//验证通过，加入 login
				temp_client_id, _ := strconv.ParseInt(claims["client_id"].(string), 10, 64)
				client_id := int32(temp_client_id)

				login_info := &LoginInfo{client: c, user_id: client_id}
				c.hub.shop_login <- login_info
				log.Printf("shop %v login success", client_id)
				c.is_login = true
				c.user_id = client_id
			} else {
				//token 验证不通过，断开连接
				if _, ok := c.hub.clients[c]; ok {
					break
				}
				fmt.Println(err)
			}
		} else if myWSMessage.Type == "message" {

			//验证链接的登录状态
			if !c.is_login {
				if _, ok := c.hub.clients[c]; ok {
					break
				}

			} else {
				//消息
				NormalMessage := &Message{}
				proto.Unmarshal(myWSMessage.Content, NormalMessage)

	/*			log.Printf("NormalMessage %v ", NormalMessage)
				log.Printf("NormalMessage.From %v ", NormalMessage.From)
				log.Printf("c.user_id %v ", c.user_id)
				log.Printf("NormalMessage.ToUserId %v ", NormalMessage.ToUserId)
				log.Printf("origin_message %v ", origin_message)*/

				//替换掉 from 的user_id
				TmpMessage := &Message{
					From:        NormalMessage.From,
					MessageType: NormalMessage.MessageType,
					FromUserId:  c.user_id,
					ToUserId:    NormalMessage.ToUserId,
					Content:     NormalMessage.Content,
				}
				TmpMessageBinary, _ := proto.Marshal(TmpMessage)
				tmp_wsMessage := &WSMessage{Type: myWSMessage.Type, Content: TmpMessageBinary}
				send_message, _ := proto.Marshal(tmp_wsMessage)

				//判断是买家消息还是卖家消息
				//判断对面的用户是否登录。
				//
				c.sqlMessage.unread_session <- TmpMessage

				if "buyer" == NormalMessage.From {
					if _, ok := c.hub.shops[NormalMessage.ToUserId]; ok {
						c.hub.shops[NormalMessage.ToUserId].send <- send_message
					}
				} else if "shoper" == NormalMessage.From {
					if _, ok := c.hub.users[NormalMessage.ToUserId]; ok {
						c.hub.users[NormalMessage.ToUserId].send <- send_message
					}
				}
			}

		}

	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			//log.Printf("message send %v ", message)

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, sqlMessage *SqlMessage, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 2560), is_login: false, user_id: 0, sqlMessage: sqlMessage}
	//todo，register 要维护一个时间，超时了直接断开。
	client.hub.register <- client
	//log.Printf("client is %v", client)
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
