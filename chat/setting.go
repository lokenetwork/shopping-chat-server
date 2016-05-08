package chat

import (

)


type setting struct {
	passport_domain string
	default_http string
	minBuyerWebsocketId uint32
	maxBuyerWebsocketId uint32
	maxShoperConnectNumber uint32
}

var Setting = setting{
	"passport.cloth.com",
	"http://",
	//Todo,the fellow will get in setting server center.
	1,
	10,
	5,
}
