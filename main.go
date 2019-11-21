package huobiapi

import (
	"github.com/bitly/go-simplejson"
	"github.com/cmdedj/huobiapi/ws"
	"github.com/cmdedj/huobiapi/client"
)

type JSON = simplejson.Json

type ParamsData = client.ParamData
type Market = ws.Market
type Listener = ws.Listener
type Client = client.Client

/// 创建WebSocket版Market客户端
func NewMarket() (*ws.Market, error) {
	return ws.NewMarket()
}

/// 创建RESTFul客户端
func NewClient(accessKeyId, accessKeySecret string) (*client.Client, error) {
	return client.NewClient(client.Endpoint, accessKeyId, accessKeySecret)
}
