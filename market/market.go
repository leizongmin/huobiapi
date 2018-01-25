package market

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"sync"

	"math"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
	"github.com/leizongmin/huobiapi/debug"
)

/// 行情的Websocket入口
var Endpoint = "wss://api.huobi.pro/ws"

type Market struct {
	inited            bool
	ws                *websocket.Conn
	wsClosed          bool
	listeners         map[string]Listener
	subscribedTopic   map[string]bool
	subscribeResultCb map[string]jsonChan
	requestResultCb   map[string]jsonChan
	lastPing          int64
	lock              sync.Mutex
}

/// 订阅事件监听器
type Listener = func(topic string, json *simplejson.Json)

/// 创建Market实例
func NewMarket() (m *Market, err error) {
	m = &Market{
		ws:                nil,
		wsClosed:          true,
		listeners:         make(map[string]Listener),
		subscribedTopic:   make(map[string]bool),
		subscribeResultCb: make(map[string]jsonChan),
		requestResultCb:   make(map[string]jsonChan),
	}
	if err := m.ReConnect(); err != nil {
		return nil, err
	}

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

		// 如果读取出错，直接关闭连接
		if err != nil {
			debug.Println(err)
			m.reconnectDelay()
			break
		}

		debug.Println("readMessage", string(msg))
		json, err := simplejson.NewJson(msg)
		if err != nil {
			debug.Println(err)
			continue
		}

		// 处理ping消息
		if ping := json.Get("ping").MustInt64(); ping > 0 {
			m.handlePing(pingData{Ping: ping})
			continue
		}

		// 处理订阅消息
		if ch := json.Get("ch").MustString(); ch != "" {
			if listener, ok := m.listeners[ch]; ok {
				debug.Println("handleSubscribe", json, msg)
				listener(ch, json)
			}
			continue
		}

		// 处理订阅成功通知
		if subbed := json.Get("subbed").MustString(); subbed != "" {
			if c, ok := m.subscribeResultCb[subbed]; ok {
				c <- json
			}
			continue
		}

		// 请求行情结果
		if rep, id := json.Get("rep").MustString(), json.Get("id").MustString(); rep != "" && id != "" {
			if c, ok := m.requestResultCb[id]; ok {
				c <- json
			}
			continue
		}

		// 处理错误消息
		if status := json.Get("status").MustString(); status == "error" {
			// 判断是否为订阅失败
			id := json.Get("id").MustString()
			if c, ok := m.subscribeResultCb[id]; ok {
				c <- json
			}
			continue
		}
	}
	debug.Println("endHandleMessageLoop")
}

/// 保持活跃
func (m *Market) keepAlive() {
	debug.Println("startKeepAlive")
	for !m.wsClosed {
		var t = getUinxMillisecond()

		// 定时主动发送ping
		debug.Println("keepAlive")
		time.Sleep(time.Second * 10)
		m.sendMessage(pingData{Ping: t})

		// 检查上次ping时间，如果超过20秒无响应，重新连接
		if math.Abs(float64(t-m.lastPing)) >= 20 {
			m.reconnectDelay()
		}
	}
	debug.Println("endKeepAlive")
}

/// 处理Ping
func (m *Market) handlePing(ping pingData) (err error) {
	debug.Println("handlePing", ping)
	m.lastPing = ping.Ping
	var pong = pongData{Pong: ping.Ping}
	err = m.sendMessage(pong)
	if err != nil {
		return err
	}
	return nil
}

/// 订阅
func (m *Market) Subscribe(topic string, listener Listener) error {
	debug.Println("subscribe", topic)
	var isNew = false
	// 如果未曾发送过订阅指令，则发送，并等待订阅操作结果，否则直接返回
	if _, ok := m.subscribedTopic[topic]; !ok {
		m.subscribeResultCb[topic] = make(jsonChan)
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
	debug.Println("unSubscribe", topic)
	// 火币网没有提供取消订阅的接口，只能删除监听器
	delete(m.listeners, topic)
}

/// 请求行情信息
func (m *Market) Request(req string) (*simplejson.Json, error) {
	var id = getRandomString(10)
	m.requestResultCb[id] = make(jsonChan)
	if err := m.sendMessage(reqData{Req: req, ID: id}); err != nil {
		return nil, err
	}
	var json = <-m.requestResultCb[id]
	delete(m.requestResultCb, id)
	// 判断是否出错
	if msg := json.Get("err-msg").MustString(); msg != "" {
		return json, fmt.Errorf(msg)
	}
	return json, nil
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
func (m *Market) Close() error {
	debug.Println("close")
	m.wsClosed = true
	return m.ws.Close()
}

/// 重新连接
func (m *Market) ReConnect() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.inited {
		debug.Println("reConnect")
	} else {
		debug.Println("connect")
	}

	if !m.wsClosed {
		if err := m.Close(); err != nil {
			return err
		}
	}
	m.wsClosed = false

	ws, _, err := websocket.DefaultDialer.Dial(Endpoint, nil)
	if err != nil {
		return err
	}
	m.ws = ws

	// 处理接收到的消息
	go m.handleMessageLoop()
	// 保持活跃
	go m.keepAlive()

	if m.inited {
		// 清理临时请求回调
		var listeners = m.listeners
		var subscribedTopic = m.subscribedTopic
		m.subscribeResultCb = make(map[string]jsonChan)
		m.requestResultCb = make(map[string]jsonChan)
		m.listeners = make(map[string]Listener)
		m.subscribedTopic = make(map[string]bool)
		// 重新订阅
		for _, topic := range getMapKeys(subscribedTopic) {
			if listener, ok := listeners[topic]; ok {
				m.Subscribe(topic, listener)
			}
		}
	}

	m.inited = true
	return nil
}

/// 稍后重新连接
func (m *Market) reconnectDelay() {
	debug.Println("reconnectDelay")
	go func() {
		time.Sleep(time.Second)
		m.Close()
		m.ReConnect()
	}()
}
