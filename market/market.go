package market

import (
	"encoding/json"
	"fmt"
	"time"

	"sync"

	"math"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
	"github.com/leizongmin/huobiapi/debug"
)

/// 行情的Websocket入口
var Endpoint = "wss://api.huobi.pro/ws"

/// 轮询时间间隔
var PollingDelay time.Duration = time.Millisecond * 100

/// Websocket未连接错误
var ConnectionClosedError = fmt.Errorf("websocket connection closed")

type wsOperation struct {
	cmd  string
	data interface{}
}

type Market struct {
	lock sync.RWMutex

	connected  bool
	userClosed bool
	destroyed  bool

	ws       *websocket.Conn
	wsOpChan chan wsOperation

	listeners         map[string]Listener
	subscribedTopic   map[string]bool
	subscribeResultCb map[string]jsonChan
	requestResultCb   map[string]jsonChan

	// 上次接收到的ping时间戳
	lastPing int64

	// 主动发送心跳的时间间隔，默认5秒
	HeartbeatInterval time.Duration
	// 接收消息超时时间，默认10秒
	ReceiveTimeout time.Duration
}

/// 订阅事件监听器
type Listener = func(topic string, json *simplejson.Json)

/// 创建Market实例
func NewMarket() (m *Market, err error) {
	m = &Market{
		HeartbeatInterval: 5 * time.Second,
		ReceiveTimeout:    10 * time.Second,
		listeners:         make(map[string]Listener),
		wsOpChan:          make(chan wsOperation),
	}
	m.initData()

	// 处理接收到的消息
	go m.handleMessageLoop()
	// 保持活跃
	go m.keepAlive()

	if err := m.connect(); err != nil {
		return nil, err
	}

	return m, nil
}

/// 重新初始化数据
func (m *Market) initData() {
	m.ws = nil
	m.connected = false
	m.subscribeResultCb = make(map[string]jsonChan)
	m.requestResultCb = make(map[string]jsonChan)
	m.subscribedTopic = make(map[string]bool)
}

/// 连接
func (m *Market) connect() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	debug.Println("connecting")
	ws, _, err := websocket.DefaultDialer.Dial(Endpoint, nil)
	if err != nil {
		return err
	}
	m.ws = ws
	m.connected = true
	m.lastPing = getUinxMillisecond()
	debug.Println("connected")

	return nil
}

/// 端口连接
func (m *Market) disconnect() (err error) {
	debug.Println("disconnecting")
	if m.ws != nil {
		m.lock.Lock()
		err = m.ws.Close()
		m.connected = false
		m.lock.Unlock()
		if err != nil {
			return err
		}
	}
	debug.Println("disconnected")
	return nil
}

/// 轮询等待时间
func (m *Market) pollingDelay() {
	time.Sleep(PollingDelay)
}

/// 稍后重新连接
func (m *Market) reconnectDelay() {
	debug.Println("reconnectDelay")
	//go func() {
	m.pollingDelay()
	if m.userClosed {
		debug.Println("user call close, skip reconnect")
		return
	}
	m.ReConnect()
	//}()
}

/// 等待连接成功再返回
func (m *Market) waitConnected() {
	if m.destroyed {
		return
	}
	for !m.connected {
		m.pollingDelay()
	}
}

/// 读取消息
func (m *Market) readMessage() (msg []byte, err error) {
	if !m.connected {
		return nil, ConnectionClosedError
	}

	// 设置接收消息超时时间
	m.lock.RLock()
	if !m.connected {
		err = ConnectionClosedError
	} else {
		err = m.ws.SetReadDeadline(time.Now().Add(m.ReceiveTimeout))
	}
	m.lock.RUnlock()
	if err != nil {
		// 判断是否为连接关闭错误
		if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
			m.reconnectDelay()
		}
		return nil, err
	}

	m.lock.RLock()
	var n int
	if !m.connected {
		err = ConnectionClosedError
	} else {
		n, msg, err = m.ws.ReadMessage()
	}
	m.lock.RUnlock()
	if err != nil {
		// 判断是否为连接关闭错误
		if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
			m.reconnectDelay()
		}
		return msg, err
	} else if n < 1 {
		return msg, nil
	} else {
		return unGzipData(msg)
	}
}

/// 发送消息
func (m *Market) sendMessage(data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	if !m.connected {
		return ConnectionClosedError
	}
	debug.Println("sendMessage", string(b))

	m.lock.Lock()
	if !m.connected {
		err = ConnectionClosedError
	} else {
		err = m.ws.WriteMessage(websocket.TextMessage, b)
	}
	m.lock.Unlock()
	if err != nil {
		return err
	}

	return nil
}

/// 处理消息循环
func (m *Market) handleMessageLoop() {
	debug.Println("startHandleMessageLoop")
	for !m.destroyed {
		m.waitConnected()
		msg, err := m.readMessage()

		// 如果读取出错，直接关闭连接
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
		if ping := json.Get("ping").MustInt64(); ping > 0 {
			go m.handlePing(pingData{Ping: ping})
			continue
		}

		// 处理pong消息
		if pong := json.Get("pong").MustInt64(); pong > 0 {
			m.lock.Lock()
			m.lastPing = pong
			m.lock.Unlock()
			continue
		}

		// 处理订阅消息
		if ch := json.Get("ch").MustString(); ch != "" {
			m.lock.RLock()
			listener, ok := m.listeners[ch]
			m.lock.RUnlock()
			if ok {
				debug.Println("handleSubscribe", json)
				go listener(ch, json)
			}
			continue
		}

		// 处理订阅成功通知
		if subbed := json.Get("subbed").MustString(); subbed != "" {
			m.lock.RLock()
			c, ok := m.subscribeResultCb[subbed]
			m.lock.RUnlock()
			if ok {
				c <- json
			}
			continue
		}

		// 请求行情结果
		if rep, id := json.Get("rep").MustString(), json.Get("id").MustString(); rep != "" && id != "" {
			m.lock.RLock()
			c, ok := m.requestResultCb[id]
			m.lock.RUnlock()
			if ok {
				c <- json
			}
			continue
		}

		// 处理错误消息
		if status := json.Get("status").MustString(); status == "error" {
			// 判断是否为订阅失败
			id := json.Get("id").MustString()
			m.lock.RLock()
			c, ok := m.subscribeResultCb[id]
			m.lock.RUnlock()
			if ok {
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
	for !m.destroyed {
		m.waitConnected()
		time.Sleep(m.HeartbeatInterval)

		// 获取当前时间戳
		var t = getUinxMillisecond()

		// 定时主动发送ping
		m.sendMessage(pingData{Ping: t})

		// 检查上次ping时间，如果超过20秒无响应，重新连接
		tr := time.Duration(math.Abs(float64(t - m.lastPing)))
		if tr >= m.HeartbeatInterval*2 {
			debug.Println("no ping max delay", tr, m.HeartbeatInterval*2, t, m.lastPing)
			m.reconnectDelay()
		}
	}
	debug.Println("endKeepAlive")
}

/// 处理Ping
func (m *Market) handlePing(ping pingData) (err error) {
	debug.Println("handlePing", ping)

	m.lock.Lock()
	m.lastPing = ping.Ping
	m.lock.Unlock()

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
		m.lock.Lock()
		m.subscribeResultCb[topic] = make(jsonChan)
		m.lock.Unlock()
		m.sendMessage(subData{ID: topic, Sub: topic})
		isNew = true
	} else {
		debug.Println("send subscribe before, reset listener only")
	}

	m.lock.Lock()
	m.listeners[topic] = listener
	m.subscribedTopic[topic] = true
	m.lock.Unlock()

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
	m.lock.Lock()
	delete(m.listeners, topic)
	m.lock.Unlock()
}

/// 请求行情信息
func (m *Market) Request(req string) (*simplejson.Json, error) {
	var id = getRandomString(10)

	m.lock.Lock()
	m.requestResultCb[id] = make(jsonChan)
	m.lock.Unlock()

	if err := m.sendMessage(reqData{Req: req, ID: id}); err != nil {
		return nil, err
	}
	var json = <-m.requestResultCb[id]

	m.lock.Lock()
	delete(m.requestResultCb, id)
	m.lock.Unlock()

	// 判断是否出错
	if msg := json.Get("err-msg").MustString(); msg != "" {
		return json, fmt.Errorf(msg)
	}
	return json, nil
}

/// 进入循环
func (m *Market) Loop() {
	debug.Println("startLoop")
	for m.connected {
		m.pollingDelay()
	}
	debug.Println("endLoop")
}

/// 重新连接
func (m *Market) ReConnect() (err error) {
	debug.Println("reconnect")

	// 如果已经连接过则先关闭连接
	if m.connected {
		debug.Println("close old connection before reconnect")
		err = m.disconnect()
	}
	if err != nil {
		return err
	}

	// 备份旧的订阅
	m.lock.RLock()
	var listeners = make(map[string]Listener)
	for k, v := range m.listeners {
		listeners[k] = v
	}
	m.lock.RUnlock()

	if err := m.connect(); err != nil {
		return err
	}

	// 恢复之前的订阅
	debug.Println("restore subscribes")
	for topic, listener := range listeners {
		m.lock.Lock()
		delete(m.subscribedTopic, topic)
		m.lock.Unlock()
		m.Subscribe(topic, listener)
	}

	return nil
}

/// 关闭连接
func (m *Market) Close() (err error) {
	debug.Println("close")
	m.userClosed = true

	err = m.disconnect()
	if err != nil {
		return err
	}

	return nil
}

/// 销毁释放资源
func (m *Market) Destroy() (err error) {
	debug.Println("destroy")

	err = m.disconnect()
	if err != nil {
		return err
	}

	m.lock.Lock()
	m.destroyed = true
	m.ws = nil
	m.listeners = nil
	m.subscribedTopic = nil
	m.subscribeResultCb = nil
	m.requestResultCb = nil
	m.lock.Unlock()

	return nil
}
