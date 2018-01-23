package market

import (
	"encoding/json"

	"github.com/leizongmin/huobiapi/client"
	"github.com/leizongmin/huobiapi/types"
)

var EndPoint = "https://api.huobi.pro/market"

type Client struct {
	base *client.Client
}

func NewClient(options client.ClientOptions) *Client {
	if options.Host == "" {
		options.Host = "api.huobi.pro"
	}
	return &Client{
		base: client.NewClient(options),
	}
}

func (c *Client) GetKLine() (*types.KLineData, error) {
	b, err := c.base.Request("GET", "/market/history/kline", nil)
	if err != nil {
		return nil, err
	}
	data := &types.KLineData{}
	err = json.Unmarshal(b, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
