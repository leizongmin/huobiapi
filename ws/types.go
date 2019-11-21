package ws

import "github.com/bitly/go-simplejson"

type pongData struct {
	Pong int64 `json:"pong"`
}

type pingData struct {
	Ping int64 `json:"ping"`
}

type subData struct {
	Sub string `json:"sub"`
	ID  string `json:"id"`
}

type reqData struct {
	Req string `json:"req"`
	ID  string `json:"id"`
}

type jsonChan = chan *simplejson.Json
