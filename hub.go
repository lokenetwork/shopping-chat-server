// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	users map[int32]*Client
	shops map[int32]*Client
	//users map[int32]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	login      chan *LoginInfo
	shop_login chan *LoginInfo

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		login:      make(chan *LoginInfo, 1),
		shop_login: make(chan *LoginInfo, 1),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		users:      make(map[int32]*Client),
		shops:      make(map[int32]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case LoginInfo := <-h.login:
			h.users[LoginInfo.user_id] = LoginInfo.client
		case LoginInfo := <-h.shop_login:
			h.shops[LoginInfo.user_id] = LoginInfo.client
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				if( client == h.users[client.user_id] ){
					delete(h.users, client.user_id)
				}else if( client == h.shops[client.user_id] ){
					delete(h.shops, client.user_id)
				}
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
