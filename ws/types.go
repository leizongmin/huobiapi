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

type AuthData struct {
	Op               string `json:"op"`
	Cid              string `json:"cid"`
	AccessKeyId      string `json:"AccessKeyId"`
	SignatureMethod  string `json:"SignatureMethod"`
	SignatureVersion string `json:"SignatureVersion"`
	Timestamp        string `json:"Timestamp"`
	Signature        string `json:"Signature"`
}

type AccountsList struct {
	Op    string `json:"op"`
	Cid   string `json:"cid"`
	Topic string `json:"topic"`
}

type jsonChan = chan *simplejson.Json

// Listener 订阅事件监听器
type Listener = func(topic string, json *simplejson.Json)
