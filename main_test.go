package huobiapi

import (
	"fmt"
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/stretchr/testify/assert"
)

func TestNewMarket(t *testing.T) {
	m, err := NewMarket()
	assert.NoError(t, err)
	go func() {
		time.Sleep(time.Second * 10)
		m.Close()
	}()
	m.Subscribe("market.eosusdt.trade.detail", func(topic string, json *simplejson.Json, raw []byte) {
		fmt.Println(topic, json, raw)
	})
	m.Loop()
}

func TestNewClient(t *testing.T) {
	c, err := NewClient("", "")
	assert.NoError(t, err)
	ret, err := c.Request("GET", "/market/history/trade", ParamsData{"symbol": "eosusdt", "size": "10"})
	assert.NoError(t, err)
	fmt.Println(ret)
}
