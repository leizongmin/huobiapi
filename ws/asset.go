package ws

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"math"
	"sync"
	"time"
	log "github.com/sirupsen/logrus"
)

// Endpoint 行情的Websocket入口
var assetEndpoint = "wss://api-aws.huobi.pro/ws/v1"

type Asset struct {
	ws *SafeWebSocket

	listeners         map[string]Listener
	listenerMutex     sync.Mutex
	subscribedTopic   map[string]bool
	subscribeResultCb map[string]jsonChan
	requestResultCb   map[string]jsonChan

	// 掉线后是否自动重连，如果用户主动执行Close()则不自动重连
	autoReconnect bool

	// 上次接收到的ping时间戳
	lastPing int64

	// 主动发送心跳的时间间隔，默认5秒
	HeartbeatInterval time.Duration
	// 接收消息超时时间，默认10秒
	ReceiveTimeout time.Duration
}

// NewMarket 创建Market实例
func NewAsset() (asset *Asset, err error) {
	asset = &Asset{
		HeartbeatInterval: 5 * time.Second,
		ReceiveTimeout:    10 * time.Second,
		ws:                nil,
		autoReconnect:     true,
		listeners:         make(map[string]Listener),
		subscribeResultCb: make(map[string]jsonChan),
		requestResultCb:   make(map[string]jsonChan),
		subscribedTopic:   make(map[string]bool),
	}

	if err := asset.connect(); err != nil {
		return nil, err
	}

	return asset, nil
}

// connect 连接
func (asset *Asset) connect() error {
	ws, err := NewSafeWebSocket(assetEndpoint)
	if err != nil {
		return err
	}
	asset.ws = ws
	asset.lastPing = getUinxMillisecond()

	asset.handleMessageLoop()
	asset.keepAlive()

	return nil
}

// reconnect 重新连接
func (asset *Asset) reconnect() error {

	time.Sleep(time.Second)

	if err := asset.connect(); err != nil {

		return err
	}

	// 重新订阅
	asset.listenerMutex.Lock()
	var listeners = make(map[string]Listener)
	for k, v := range asset.listeners {
		listeners[k] = v
	}
	asset.listenerMutex.Unlock()

	for topic, listener := range listeners {
		delete(asset.subscribedTopic, topic)
		asset.Subscribe(topic, listener)
	}
	return nil
}

// sendMessage 发送消息
func (asset *Asset) sendMessage(data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	log.Info("send message: ", string(b))
	asset.ws.Send(b)
	return nil
}

// handleMessageLoop 处理消息循环
func (asset *Asset) handleMessageLoop() {
	asset.ws.Listen(func(buf []byte) {
		msg, err := unGzipData(buf)

		if err != nil {

			return
		}
		json, err := simplejson.NewJson(msg)
		if err != nil {

			return
		}

		log.Info(json)

		op := json.Get("op").MustString()

		// 处理ping
		if op == "ping"{
			ts := json.Get("ts").MustInt64()
			err := asset.handlePing(pingData{
				Op: "ping",
				Ts: ts,
			})
			if err != nil {
				log.Error(err)
			}

		}


		// 处理pong消息
		if pong := json.Get("pong").MustInt64(); pong > 0 {
			asset.lastPing = pong
			return
		}

		// 处理订阅消息
		if ch := json.Get("ch").MustString(); ch != "" {
			asset.listenerMutex.Lock()
			listener, ok := asset.listeners[ch]
			asset.listenerMutex.Unlock()
			if ok {

				listener(ch, json)
			}
			return
		}

		// 处理订阅成功通知
		if subbed := json.Get("subbed").MustString(); subbed != "" {
			c, ok := asset.subscribeResultCb[subbed]
			if ok {
				c <- json
			}
			return
		}

		// 请求行情结果
		if rep, id := json.Get("rep").MustString(), json.Get("id").MustString(); rep != "" && id != "" {
			c, ok := asset.requestResultCb[id]
			if ok {
				c <- json
			}
			return
		}

		// 处理错误消息
		if status := json.Get("status").MustString(); status == "error" {
			// 判断是否为订阅失败
			id := json.Get("id").MustString()
			c, ok := asset.subscribeResultCb[id]
			if ok {
				c <- json
			}
			return
		}
	})
}

// keepAlive 保持活跃
func (asset *Asset) keepAlive() {
	asset.ws.KeepAlive(asset.HeartbeatInterval, func() {
		var t = getUinxMillisecond()
		// asset.sendMessage(pingData{Ping: t})

		// 检查上次ping时间，如果超过20秒无响应，重新连接
		tr := time.Duration(math.Abs(float64(t - asset.lastPing)))
		if tr >= asset.HeartbeatInterval*2 {

			if asset.autoReconnect {
				err := asset.reconnect()
				if err != nil {

				}
			}
		}
	})
}

// handlePing 处理Ping
func (asset *Asset) handlePing(ping pingData) (err error) {

	asset.lastPing = ping.Ts
	var pong = pongData{
		Op: "pong",
		Ts: ping.Ts,
	}
	err = asset.sendMessage(pong)
	if err != nil {
		return err
	}
	return nil
}

// Subscribe 订阅
func (asset *Asset) Subscribe(topic string, listener Listener) error {


	var isNew = false

	// 如果未曾发送过订阅指令，则发送，并等待订阅操作结果，否则直接返回
	if _, ok := asset.subscribedTopic[topic]; !ok {
		asset.subscribeResultCb[topic] = make(jsonChan)
		asset.sendMessage(subData{ID: topic, Sub: topic})
		isNew = true
	} else {

	}

	asset.listenerMutex.Lock()
	asset.listeners[topic] = listener
	asset.listenerMutex.Unlock()
	asset.subscribedTopic[topic] = true

	if isNew {
		var json = <-asset.subscribeResultCb[topic]
		// 判断订阅结果，如果出错则返回出错信息
		if msg, err := json.Get("err-msg").String(); err == nil {
			return fmt.Errorf(msg)
		}
	}
	return nil
}

// Unsubscribe 取消订阅
func (asset *Asset) Unsubscribe(topic string) {

	asset.listenerMutex.Lock()
	// 火币网没有提供取消订阅的接口，只能删除监听器
	delete(asset.listeners, topic)
	asset.listenerMutex.Unlock()
}

// Request 请求行情信息
func (asset *Asset) Request(req string) (*simplejson.Json, error) {
	var id = getRandomString(10)
	asset.requestResultCb[id] = make(jsonChan)

	if err := asset.sendMessage(reqData{Req: req, ID: id}); err != nil {
		return nil, err
	}
	var json = <-asset.requestResultCb[id]

	delete(asset.requestResultCb, id)

	// 判断是否出错
	if msg := json.Get("err-msg").MustString(); msg != "" {
		return json, fmt.Errorf(msg)
	}
	return json, nil
}

// Loop 进入循环
func (asset *Asset) Loop() {

	for {
		err := asset.ws.Loop()
		if err != nil {

			if err == SafeWebSocketDestroyError {
				break
			} else if asset.autoReconnect {
				asset.reconnect()
			} else {
				break
			}
		}
	}

}

// ReConnect 重新连接
func (asset *Asset) ReConnect() (err error) {

	asset.autoReconnect = true
	if err = asset.ws.Destroy(); err != nil {
		return err
	}
	return asset.reconnect()
}

// Close 关闭连接
func (asset *Asset) Close() error {

	asset.autoReconnect = false
	if err := asset.ws.Destroy(); err != nil {
		return err
	}
	return nil
}


func (asset * Asset) Auth(accessKeyId, accessKeySecret string) error {
	params := make(map[string]string)

	params["AccessKeyId"] = accessKeyId
	params["SignatureMethod"] = "HmacSHA256"
	params["SignatureVersion"] = "2"
	params["Timestamp"] = time.Now().Format("2006-01-02T15:04:05")

	authData := AuthData{
		Op:              "auth",
		Cid:              "xxxxxxxxxxxxx",
		AccessKeyId:      params["AccessKeyId"],
		SignatureMethod:  params["SignatureMethod"],
		SignatureVersion: params["SignatureVersion"],
		Timestamp:        params["Timestamp"],
		Signature:        GenSignature(params, accessKeySecret),
	}

	err := asset.sendMessage(authData)
	return err
}
