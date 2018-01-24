package market

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	//"github.com/bitly/go-simplejson"
	"github.com/bitly/go-simplejson"
)

func TestMarket(t *testing.T) {
	m, err := NewMarket()
	assert.NoError(t, err)

	// 订阅
	err = m.Subscribe("market.eosusdt.kline.1min", func(topic string, json *simplejson.Json, raw []byte) {
		fmt.Println(topic, json, raw)
	})
	assert.NoError(t, err)
	err = m.Subscribe("market.eosusdt.trade.detail", func(topic string, json *simplejson.Json, raw []byte) {
		fmt.Println(topic, json, raw)
	})
	assert.NoError(t, err)

	// 请求
	rep, err := m.Request("market.eosusdt.detail")
	assert.NoError(t, err)
	fmt.Println(rep)

	// 阻塞事件循环
	fmt.Println(m)
	go func() {
		time.Sleep(time.Second * 10)
		m.Close()
	}()
	m.Loop()

	// 重新连接
	m.ReConnect()
	go func() {
		time.Sleep(time.Second * 10)
		m.Close()
	}()
	m.Loop()
}
