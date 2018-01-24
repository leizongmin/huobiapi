package market

import "github.com/bitly/go-simplejson"

type pongData struct {
	Pong int `json:"pong"`
}

type pingData struct {
	Ping int `json:"ping"`
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
