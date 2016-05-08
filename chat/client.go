package chat

import (
	"fmt"
	"io"
	"log"
	"golang.org/x/net/websocket"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const channelBufSize = 100


// Chat client.
type BuyerClient struct {
	websocketId uint32
	uid         uint32
	ws          *websocket.Conn
	server      *Server
	ch          chan *CommandMessage
	doneCh      chan bool
	loginStatus bool
}

type ShoperClient struct {
	websocketId uint32
	uid         uint32
	ws          *websocket.Conn
	server      *Server
	ch          chan *CommandMessage
	doneCh      chan bool
	loginStatus bool
}

// Create new buyer chat client.
func NewBuyerClient(ws *websocket.Conn, server *Server) (bool, *BuyerClient) {

	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	connect_success := false

	var client_id uint32 = get_the_buyer_can_use_websocket_id(server)

	if ( client_id != 0 ) {
		connect_success = true
	} else {
		ws.Close()
	}

	ch := make(chan *CommandMessage, channelBufSize)
	doneCh := make(chan bool)
	loginStatus := false

	return connect_success, &BuyerClient{client_id, 0, ws, server, ch, doneCh, loginStatus}
}

func NewShoperClient(ws *websocket.Conn, server *Server) (bool, *ShoperClient) {

	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	websocket_id := get_the_shoper_can_use_websocket_id(server)
	connect_success := false
	if ( uint32(len(server.shoperClients)) < Setting.maxShoperConnectNumber ) {
		connect_success = true
	} else {
		ws.Close()
	}

	ch := make(chan *CommandMessage, channelBufSize)
	doneCh := make(chan bool)
	loginStatus := false
	return connect_success, &ShoperClient{websocket_id, 15, ws, server, ch, doneCh, loginStatus}
}

func (c *BuyerClient) Conn() *websocket.Conn {
	return c.ws
}

func (c *BuyerClient) Write(msg *CommandMessage) {
	select {
	case c.ch <- msg:
	default:
		c.server.DelBuyerClient(c)
		err := fmt.Errorf("client %d is disconnected.", c.websocketId)
		c.server.Err(err)
	}
}

func (c *ShoperClient) Write(msg *CommandMessage) {
	select {
	case c.ch <- msg:
	default:
		c.server.DelShoperClient(c)
		err := fmt.Errorf("shop client %d is disconnected.", c.websocketId)
		c.server.Err(err)
	}
}

func (c *BuyerClient) Done() {
	log.Println("Done")
	c.doneCh <- true
}

// Listen Write and Read request via chanel
func (c *BuyerClient) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen Write and Read request via chanel
func (shopClient *ShoperClient) Listen() {
	go shopClient.listenWrite()
	shopClient.listenRead()
}

// Listen write request via chanel
func (c *ShoperClient) listenWrite() {
	log.Println("Listening write to ShoperClient")
	for {
		select {

		// send message to the client
		case msg := <-c.ch:
			websocket.JSON.Send(c.ws, msg)
		// receive done request
		case <-c.doneCh:
			c.server.DelShoperClient(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (c *ShoperClient) listenRead() {
	log.Println("Listening read from ShoperClient")
	for {
		select {

		// receive done request
		case <-c.doneCh:

			c.server.DelShoperClient(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
		//get the command
			var msg CommandMessage

			err := websocket.JSON.Receive(c.ws, &msg)

		//EOF means close.
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				if msg.Command == "set_login_info" {
					loginStatus, user_id := set_login_info(msg.Content)
					c.loginStatus = loginStatus
					c.uid = user_id
				} else if msg.Command == "send_shoper_message" {
					send_shoper_message(c, msg.Content)
				}
			}
		}
	}
}


// Listen write request via chanel
func (c *BuyerClient) listenWrite() {
	log.Println("Listening write to BuyerClient")
	for {
		select {

		// send message to the client
		case msg := <-c.ch:
			log.Println("Send:", msg)
			websocket.JSON.Send(c.ws, msg)

		// receive done request
		case <-c.doneCh:

			c.server.DelBuyerClient(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (c *BuyerClient) listenRead() {
	log.Println("Listening read from BuyerClient")
	for {
		select {

		// receive done request
		case <-c.doneCh:
			c.server.DelBuyerClient(c)
			c.doneCh <- true // for listenWrite method
			return
		// read data from websocket connection
		default:
		//get command
			var msg CommandMessage

			err := websocket.JSON.Receive(c.ws, &msg)

		//EOF means close.
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				if msg.Command == "set_login_info" {
					loginStatus, user_id := set_login_info(msg.Content)
					c.loginStatus = loginStatus
					c.uid = user_id
				} else if msg.Command == "send_buyer_message" {
					log.Println("buyer message is ", msg)
					send_buyer_message(c, msg.Content)
				}
			}


		}
	}
}

func send_buyer_message(buyerClient *BuyerClient, content *json.RawMessage) {
	var fromBuyerMsg FromBuyerMessage
	err := json.Unmarshal(*content, &fromBuyerMsg)
	if err != nil {
		buyerClient.server.Err(err)
	}

	var toShoperMsg ToShoperMessage
	toShoperMsg.FromBuyerUserId = buyerClient.uid
	toShoperMsg.FromBuyerWebsocketId = buyerClient.websocketId
	toShoperMsg.Message = fromBuyerMsg.Message

	var toShoperCmdMsg CommandMessage
	toShoperCmdMsg.Command = "receive_message"

	var toShoperRawMsg json.RawMessage
	toShoperRawMsg, err = json.Marshal(&toShoperMsg)
	if err != nil {
		buyerClient.server.Err(err)
	}
	toShoperCmdMsg.Content = &toShoperRawMsg

	//toShoperCmdMsg.Content = json.RawMessage(toShoperMsg)


	//Check the service client socket is int this computer or not.If in, end the message
	go func() {
		for _, _shopClient := range buyerClient.server.shoperClients {
			if fromBuyerMsg.ToCustomerServiceUserId == _shopClient.uid {
				log.Println("send message success")
				_shopClient.Write(&toShoperCmdMsg)
			}
		}
	}()

	//synchronous buyer message in all device
	if buyerClient.loginStatus == true {
		//this computer loop
		go func() {
			for _, otherBuyerClient := range buyerClient.server.buyerClients {
				if buyerClient.uid == otherBuyerClient.uid {
					otherBuyerClient.Write(&toShoperCmdMsg)
				}
			}
		}()

		//Todo, send the message to other server.synchronous buyer message
	}


	//todo,synchronous shopper message in all device

}

func send_shoper_message(shoperClient *ShoperClient, content *json.RawMessage) {
	var fromShoperMsg FromShoperMessage

	err := json.Unmarshal(*content, &fromShoperMsg)

	if err != nil {
		shoperClient.server.Err(err)
	}


	var toBuyerMsg ToBuyerMessage
	toBuyerMsg.FromCustomerServiceUserId = shoperClient.uid
	toBuyerMsg.Message = fromShoperMsg.Message

	var toBuyerCmdMsg CommandMessage
	toBuyerCmdMsg.Command = "receive_message"

	var toBuyerRawMsg json.RawMessage
	toBuyerRawMsg, err = json.Marshal(&toBuyerMsg)
	if err != nil {
		shoperClient.server.Err(err)
	}
	toBuyerCmdMsg.Content = &toBuyerRawMsg


	//check is the buery is login or not
	if fromShoperMsg.toBuyerUserId > 0 {
		//synchronous buyer message in all device
		go func() {
			log.Println(shoperClient.server.buyerClients)
			for _, buyerClient := range shoperClient.server.buyerClients {
				if shoperClient.server.buyerClients[fromShoperMsg.ToBuyerWebsokcetId].uid == buyerClient.uid && shoperClient.server.buyerClients[fromShoperMsg.ToBuyerWebsokcetId].websocketId != buyerClient.websocketId {
					log.Println("send shoper message success")
					buyerClient.Write(&toBuyerCmdMsg)
				}
			}
		}()
		//send to other server,synchronous buyer message
	} else {
		//If socket connect in this computer
		if fromShoperMsg.ToBuyerWebsokcetId >= Setting.minBuyerWebsocketId && fromShoperMsg.ToBuyerWebsokcetId <= Setting.maxBuyerWebsocketId {

			if ( shoperClient.server.buyerClients[fromShoperMsg.ToBuyerWebsokcetId] != nil ) {
				shoperClient.server.buyerClients[fromShoperMsg.ToBuyerWebsokcetId].Write(&toBuyerCmdMsg)
			}
		} else {
			//todo,socket connect not in this computer
		}
	}

}

func set_login_info(content *json.RawMessage) (bool, uint32) {

	var loginMsg LoginMessage
	err := json.Unmarshal(*content, &loginMsg)
	if err != nil {
		return false, 0
	}

	get_login_info_ur := Setting.default_http + Setting.passport_domain + "/User/echoJsonInfo?passport_login_key=" + (loginMsg.LoginToken)
	resp, err := http.Get(get_login_info_ur)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	passport_info := map[string]interface{}{}
	json.Unmarshal(body, &passport_info)
	var passport_user_info = passport_info["passport_user_info"].(map[string]interface{})
	var passport_login_status float64 = passport_info["passport_login_status"].(float64)
	loginStatus := false
	var user_id uint32 = 0
	if passport_login_status == 1 {
		loginStatus = true
		user_id = uint32(passport_user_info["user_id"].(float64))
	}
	return loginStatus, user_id
}


/*
 * 0 stand that,connect number is fill will.
 */
func get_the_buyer_can_use_websocket_id(server *Server) uint32 {
	var client_index uint32
	//获取可用的client number
	for client_index = Setting.minBuyerWebsocketId; client_index <= Setting.minBuyerWebsocketId; client_index++ {
		log.Println("client_index is", client_index)
		if server.buyerClients[client_index] == nil {
			break
		}
	}
	if ( server.buyerClients[client_index] == nil ) {
		return client_index;
	} else {
		return 0
	}
}

/*
 * 0 stand that,connect number is fill will.
 */
func get_the_shoper_can_use_websocket_id(server *Server) uint32 {
	var client_index uint32
	//获取可用的client number
	for client_index = 0; client_index <= Setting.maxShoperConnectNumber; client_index++ {
		if server.shoperClients[client_index] == nil {
			break
		}
	}
	return client_index;

}


