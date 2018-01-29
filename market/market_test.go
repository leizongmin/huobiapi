package market

import (
	"fmt"
	"testing"
	"time"

	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/stretchr/testify/assert"
)

func TestNewMarket(t *testing.T) {
	m, err := NewMarket()
	assert.NoError(t, err)

	// 订阅
	err = m.Subscribe("market.eosusdt.kline.1min", func(topic string, json *simplejson.Json) {
		fmt.Println(topic, json)
	})
	assert.NoError(t, err)
	err = m.Subscribe("market.eosusdt.trade.detail", func(topic string, json *simplejson.Json) {
		fmt.Println(topic, json)
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

	fmt.Println(strings.Repeat("-------------------\n", 10))

	// 重新连接
	m.ReConnect()
	go func() {
		time.Sleep(time.Second * 12)
		m.Close()
	}()
	go func() {
		for {
			time.Sleep(time.Second * 2)
			rep, err := m.Request("market.eosusdt.detail")
			fmt.Println(err, rep)
		}
	}()
	m.Loop()

	fmt.Println(m)
}

func TestMarketAlive(t *testing.T) {
	m, err := NewMarket()
	assert.NoError(t, err)
	err = m.Subscribe("market.eosusdt.kline.1min", func(topic string, json *simplejson.Json) {
		fmt.Println(topic, json)
	})
	assert.NoError(t, err)
	go func() {
		time.Sleep(time.Minute * 10)
		m.Close()
	}()
	m.Loop()
}
