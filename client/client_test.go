package client

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/bitly/go-simplejson"
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
	body, err := SendRequest(sign, "GET", "api.huobi.pro", "/market/history/kline", ParamData{
		"period": "1day",
		"size":   "200",
		"symbol": "btcusdt",
	})
	assert.NoError(t, err)
	fmt.Println(string(body))
	js, err := simplejson.NewJson(body)
	assert.NoError(t, err)
	if status, err := js.Get("status").String(); err != nil {
		assert.NoError(t, err)
	} else {
		assert.Equal(t, "ok", status)
	}
}
