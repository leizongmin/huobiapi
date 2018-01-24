package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	sign := NewSign("e2xxxxxx-99xxxxxx-84xxxxxx-7xxxx", "b0xxxxxx-c6xxxxxx-94xxxxxx-dxxxx")
	ret, err := sign.Get("GET", "api.huobi.pro", "/v1/order/orders", "2017-05-11T15:19:30", map[string]string{
		"order-id": "1234567890",
	})
	assert.NoError(t, err)
	fmt.Println(ret)
	assert.Equal(t, "Nmd8AU8uAe0mkFpxNbiava0aeZzBEtYjCdie1ZYZjoM=", ret)
}

func TestSendRequest(t *testing.T) {
	sign := NewSign("e2xxxxxx-99xxxxxx-84xxxxxx-7xxxx", "b0xxxxxx-c6xxxxxx-94xxxxxx-dxxxx")
	json, err := SendRequest(sign, "GET", "https", "api.huobi.pro", "/market/history/kline", ParamData{
		"period": "1day",
		"size":   "200",
		"symbol": "btcusdt",
	})
	assert.NoError(t, err)
	fmt.Println(json)
}

func TestClient_Request(t *testing.T) {
	client, err := NewClient(MarketEndpoint, "", "")
	assert.NoError(t, err)
	json, err := client.Request("GET", "/trade", ParamData{"symbol": "eosusdt"})
	assert.NoError(t, err)
	fmt.Println(json)

	client, err = NewClient(Endpoint, "", "")
	assert.NoError(t, err)
	json, err = client.Request("GET", "/market/history/trade", ParamData{"symbol": "eosusdt", "size": "10"})
	assert.NoError(t, err)
	fmt.Println(json)
}
