package chat

import "encoding/json"

type CommandMessage struct {
	Command string `json:"command"`
	Content *json.RawMessage `json:"content"`
}

type LoginMessage struct {
	LoginToken string `json:"loginToken"`
}

type NormalMessage struct {
	ToWebsocketId uint32 `json:"toWebsocketId"`
	Message string `json:"message"`
}

type FromBuyerMessage struct {
	ToCustomerServiceUserId uint32 `json:"toCustomerServiceUserId"`
	Message string `json:"message"`
}

type FromShoperMessage struct {
	toBuyerUserId uint32 `json:"toBuyerWebsokcetId"`
	ToBuyerWebsokcetId uint32 `json:"toBuyerWebsokcetId"`
	Message string `json:"message"`
}

type ToShoperMessage struct {
	FromBuyerWebsocketId uint32 `json:"fromBuyerWebsocketId"`
	FromBuyerUserId uint32 `json:"fromBuyerUserId"`
	Message string `json:"message"`
}

type ToBuyerMessage struct {
	FromCustomerServiceUserId uint32 `json:"fromCustomerServiceUserId"`
	Message string `json:"message"`
}



func (self *CommandMessage) String() string {
	return self.Command + " to "
}
