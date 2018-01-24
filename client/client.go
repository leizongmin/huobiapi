package client

import (
	"net/url"

	"github.com/bitly/go-simplejson"
)

type Client struct {
	Sign       *Sign
	host       string
	pathPrefix string
	scheme     string
}

/// 行情API
const MarketEndpoint = "https://api.huobi.pro/market"

/// 交易API
const TradeEndpoint = "https://api.huobi.pro/v1"

/// 全局API
const Endpoint = "https://api.huobi.pro"

/// 创建新客户端
func NewClient(endpoint, accessKeyId, accessKeySecret string) (*Client, error) {
	urlInfo, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Sign:       NewSign(accessKeyId, accessKeySecret),
		host:       urlInfo.Host,
		pathPrefix: urlInfo.Path,
		scheme:     urlInfo.Scheme,
	}
	if client.pathPrefix == "/" {
		client.pathPrefix = ""
	}

	return client, nil
}

/// 发送请求
func (c *Client) Request(method, path string, data ParamData) (*simplejson.Json, error) {
	return SendRequest(c.Sign, method, c.scheme, c.host, c.pathPrefix+path, data)
}
