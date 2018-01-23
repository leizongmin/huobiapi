package market_ws

import (
	"compress/gzip"
	"bytes"
	"io/ioutil"
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"time"
)

/// 行情的Websocket入口
var Endpoint = "wss://api.huobi.pro/ws"

type Market struct {
	ws *websocket.Conn
}

func NewMarket() (m *Market, err error) {
	m = &Market{}
	ws, _, err := websocket.DefaultDialer.Dial(Endpoint, nil)
	if err != nil {
		return nil, err
	}
	m.ws = ws

	m.sendMessage(subData{ID: "xxx", Sub: "market.btcusdt.kline.1min"})
	go m.handleMessageLoop()
	time.Sleep(time.Second * 10)
	m.sendMessage(pingData{Ping: int(time.Now().UnixNano()/1000000)})
	time.Sleep(time.Second * 20)

	return m, nil
}

func (m *Market) readMessage() (msg []byte, err error) {
	if n, buf, err := m.ws.ReadMessage(); err != nil {
		return msg, err
	} else if n < 1 {
		return msg, nil
	} else {
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

func (m *Market) sendMessage(data interface{}) error {
	if b, err := json.Marshal(data); err != nil {
		return err
	} else {
		log.Println(string(b))
		if err := m.ws.WriteMessage(websocket.TextMessage, b); err != nil {
			return err
		}
	}
	return nil
}

func (m *Market) handleMessageLoop() {
	for {
		msg, err := m.readMessage()
		if err != nil {
			log.Println(err)
		}
		log.Println(string(msg))
		json, err := simplejson.NewJson(msg)
		if err != nil {
			log.Println(err)
		}
		if ping, ok := json.CheckGet("ping"); ok {
			if v, err := ping.Int(); err != nil {
				log.Println(err)
			} else {
				m.handlePing(pingData{Ping: v})
			}
		}
	}
}

func (m *Market) handlePing(ping pingData) (err error) {
	var pong = pongData{Pong: ping.Ping}
	err = m.sendMessage(pong)
	if err != nil {
		return err
	}
	return nil
}

func (m *Market) Subscribe() {

}

func (m *Market) Unsubscribe() {

}
