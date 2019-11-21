package huobiapi

import (
	"github.com/bitly/go-simplejson"
	"github.com/leizongmin/huobiapi/market"
	"huobiapi/client"
	"huobiapi/ws"
)

type JSON = simplejson.Json

type ParamsData = client.ParamData
type Market = ws.Market
type Listener = ws.Listener
type Client = client.Client

/// 创建WebSocket版Market客户端
func NewMarket() (*market.Market, error) {
	return market.NewMarket()
}

/// 创建RESTFul客户端
func NewClient(accessKeyId, accessKeySecret string) (*client.Client, error) {
	return client.NewClient(client.Endpoint, accessKeyId, accessKeySecret)
}
