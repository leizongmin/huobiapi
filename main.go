package huobiapi

import (
	"github.com/bitly/go-simplejson"
	"github.com/leizongmin/huobiapi/client"
	"github.com/leizongmin/huobiapi/market"
)

type ParamsData = client.ParamData

type JSON = simplejson.Json

/// 创建WebSocket版Market客户端
func NewMarket() (*market.Market, error) {
	return market.NewMarket()
}

/// 创建RESTFul客户端
func NewClient(accessKeyId, accessKeySecret string) (*client.Client, error) {
	return client.NewClient(client.Endpoint, accessKeyId, accessKeySecret)
}
