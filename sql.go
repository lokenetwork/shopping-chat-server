// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	_ "github.com/go-sql-driver/mysql"
)

type SqlMessage struct {
	unread_session chan *Message
}

func newSqlMessage() *SqlMessage {
	return &SqlMessage{
		unread_session: make(chan *Message,100),
	}
}

func unread_session(message *Message) {
	//log.Printf("message.FromUserId is : %v", message.FromUserId)
	//log.Printf("message.ToUserId is : %v", message.ToUserId)
	//查询出会话是否存在
	//把消息写入数据库，可以查看历史消息
	var client_id int32
	var shoper_id int32
	if "buyer" == message.From {
		client_id = message.FromUserId
		shoper_id = message.ToUserId
	} else if "shoper" == message.From {
		client_id = message.ToUserId
		shoper_id = message.FromUserId
	}

	rows, _ := db.Query("SELECT session_id FROM c_new_chat_session WHERE client_id = ? AND shop_id = ?", client_id, shoper_id)
	has_session := 0
	for rows.Next() {
		has_session++
	}
	//会话不存在，插入会话
	if 0 == has_session {
		stmt2, _ := db.Prepare("INSERT INTO  c_new_chat_session(client_id,shop_id) VALUES(  ?, ?)")
		stmt2.Exec(client_id, shoper_id)
	} else {
		if "buyer" == message.From {
			stmt2, _ := db.Prepare("UPDATE c_new_chat_session SET shop_read = 0 WHERE client_id = ? AND shop_id = ?")
			stmt2.Exec(client_id, shoper_id)
		} else if "shoper" == message.From {
			stmt2, _ := db.Prepare("UPDATE c_new_chat_session SET client_read = 0 WHERE client_id = ? AND shop_id = ?")
			stmt2.Exec(client_id, shoper_id)
		}
	}
	//写入历史消息
	stmt3, _ := db.Prepare("INSERT INTO  c_new_chat_history(message_type,sender_type,shoper_id,client_id,content) VALUES(  ?, ?,?,?,?)")
	stmt3.Exec(message.MessageType, message.From, shoper_id, client_id, message.Content)
}

func (sqlMessage *SqlMessage) run() {
	for {
		select {
		case sql_message := <-sqlMessage.unread_session:
			unread_session(sql_message)
		}
	}
}
