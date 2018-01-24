package market_ws

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gorilla/websocket"
	"github.com/leizongmin/huobiapi/debug"
	"github.com/bitly/go-simplejson"
)

/// 行情的Websocket入口
var Endpoint = "wss://api.huobi.pro/ws"

type Market struct {
	ws                *websocket.Conn
	wsClosed          bool
	listeners         map[string]Listener
	subscribedTopic   map[string]bool
	subscribeResultCb map[string]chan *simplejson.Json
}

/// 订阅事件监听器
type Listener = func(topic string, json *simplejson.Json, raw []byte)

/// 创建Market实例
func NewMarket() (m *Market, err error) {
	m = &Market{
		listeners:         make(map[string]Listener),
		subscribedTopic:   make(map[string]bool),
		subscribeResultCb: make(map[string]chan *simplejson.Json),
	}

	ws, _, err := websocket.DefaultDialer.Dial(Endpoint, nil)
	if err != nil {
		return nil, err
	}
	m.ws = ws

	go m.handleMessageLoop()

	return m, nil
}

/// 读取消息
func (m *Market) readMessage() (msg []byte, err error) {
	if n, buf, err := m.ws.ReadMessage(); err != nil {
		// 判断是否为连接关闭错误
		if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
			m.wsClosed = true
			m.ws.Close()
		}
		return msg, err
	} else if n < 1 {
		return msg, nil
	} else {
		// 接受到的数据要gzip解压
		if r, err := gzip.NewReader(bytes.NewBuffer(buf)); err != nil {
			return msg, nil
		} else {
			if msg, err := ioutil.ReadAll(r); err != nil {
				return msg, nil
			} else {
				return msg, nil
			}
		}
	}
}

/// 发送消息
func (m *Market) sendMessage(data interface{}) error {
	if b, err := json.Marshal(data); err != nil {
		return err
	} else {
		debug.Println("sendMessage", string(b))
		if err := m.ws.WriteMessage(websocket.TextMessage, b); err != nil {
			return err
		}
	}
	return nil
}

/// 处理消息循环
func (m *Market) handleMessageLoop() {
	debug.Println("startHandleMessageLoop")
	for !m.wsClosed {
		msg, err := m.readMessage()
		if err != nil {
			debug.Println(err)
			continue
		}
		debug.Println("readMessage", string(msg))
		json, err := simplejson.NewJson(msg)
		if err != nil {
			debug.Println(err)
			continue
		}

		// 处理ping消息
		if ping, ok := json.CheckGet("ping"); ok {
			if v, err := ping.Int(); err != nil {
				debug.Println(err)
			} else {
				m.handlePing(pingData{Ping: v})
			}
			continue
		}

		// 处理订阅消息
		if ch, ok := json.CheckGet("ch"); ok {
			if v, err := ch.String(); err != nil {
				debug.Println(err)
			} else {
				if listener, ok := m.listeners[v]; ok {
					debug.Println("handleSubscribe", json, msg)
					listener(v, json, msg)
				}
			}
			continue
		}

		// 处理订阅成功通知
		if subbed, ok := json.CheckGet("subbed"); ok {
			if v, err := subbed.String(); err != nil {
				debug.Println(err)
			} else {
				if c, ok := m.subscribeResultCb[v]; ok {
					c <- json
				}
			}
			continue
		}

		// 处理错误消息
		if status, ok := json.CheckGet("status"); ok {
			if v, err := status.String(); err != nil {
				debug.Println(err)
			} else {
				if v == "error" {
					// 判断是否为订阅失败
					if id, err := json.Get("id").String(); err == nil {
						if c, ok := m.subscribeResultCb[id]; ok {
							c <- json
						}
					}
				}
			}
			continue
		}
	}
	debug.Println("endHandleMessageLoop")
}

/// 处理Ping
func (m *Market) handlePing(ping pingData) (err error) {
	debug.Println("handlePing", ping)
	var pong = pongData{Pong: ping.Ping}
	err = m.sendMessage(pong)
	if err != nil {
		return err
	}
	return nil
}

/// 订阅
func (m *Market) Subscribe(topic string, listener Listener) error {
	var isNew = false
	// 如果未曾发送过订阅指令，则发送，并等待订阅操作结果，否则直接返回
	if _, ok := m.subscribedTopic[topic]; !ok {
		m.subscribeResultCb[topic] = make(chan *simplejson.Json)
		m.sendMessage(subData{ID: topic, Sub: topic})
		isNew = true
	}

	m.listeners[topic] = listener
	m.subscribedTopic[topic] = true

	if isNew {
		var json = <-m.subscribeResultCb[topic]
		// 判断订阅结果，如果出错则返回出错信息
		if msg, err := json.Get("err-msg").String(); err == nil {
			return fmt.Errorf(msg)
		}
	}
	return nil
}

/// 取消订阅
func (m *Market) Unsubscribe(topic string) {
	// 火币网没有提供取消订阅的接口，只能删除监听器
	delete(m.listeners, topic)
}

/// 进入循环
func (m *Market) Loop() {
	debug.Println("startLoop")
	for !m.wsClosed {
		time.Sleep(time.Second)
	}
	debug.Println("endLoop")
}

/// 关闭
func (m *Market) Close() {
	debug.Println("close")
	m.wsClosed = true
	m.ws.Close()
}
