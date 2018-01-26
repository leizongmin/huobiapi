package market

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

var SafeWebSocketDestroyError = fmt.Errorf("connection destroy by user")

type SafeWebSocket struct {
	ws               *websocket.Conn
	listener         SafeWebSocketMessageListener
	aliveHandler     SafeWebSocketAliveHandler
	aliveInterval    time.Duration
	sendQueue        chan []byte
	lastError        error
	runningTaskSend  bool
	runningTaskRead  bool
	runningTaskAlive bool
}

type SafeWebSocketMessageListener = func(b []byte)
type SafeWebSocketAliveHandler = func()

func NewSafeWebSocket(endpoint string) (*SafeWebSocket, error) {
	ws, _, err := websocket.DefaultDialer.Dial(Endpoint, nil)
	if err != nil {
		return nil, err
	}
	s := &SafeWebSocket{ws: ws, sendQueue: make(chan []byte, 1000), aliveInterval: time.Second * 60}

	go func() {
		s.runningTaskSend = true
		for s.lastError == nil {
			b := <-s.sendQueue
			err := s.ws.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				s.lastError = err
				break
			}
		}
		s.runningTaskSend = false
	}()

	go func() {
		s.runningTaskRead = true
		for s.lastError == nil {
			_, b, err := s.ws.ReadMessage()
			if err != nil {
				s.lastError = err
				break
			}
			s.listener(b)
		}
		s.runningTaskRead = false
	}()

	go func() {
		s.runningTaskAlive = true
		for s.lastError == nil {
			if s.aliveHandler != nil {
				s.aliveHandler()
			}
			time.Sleep(s.aliveInterval)
		}
		s.runningTaskAlive = false
	}()

	return s, nil
}

func (s *SafeWebSocket) Listen(h SafeWebSocketMessageListener) {
	s.listener = h
}

func (s *SafeWebSocket) Send(b []byte) {
	s.sendQueue <- b
}

func (s *SafeWebSocket) KeepAlive(v time.Duration, h SafeWebSocketAliveHandler) {
	s.aliveInterval = v
	s.aliveHandler = h
}

func (s *SafeWebSocket) Destroy() error {
	s.lastError = SafeWebSocketDestroyError
	for !s.runningTaskRead && !s.runningTaskSend && !s.runningTaskAlive {
		time.Sleep(time.Millisecond * 100)
	}
	err := s.ws.Close()
	s.ws = nil
	s.listener = nil
	s.aliveHandler = nil
	s.sendQueue = nil
	return err
}

func (s *SafeWebSocket) Loop() error {
	for s.lastError == nil {
		time.Sleep(time.Millisecond * 100)
	}
	return s.lastError
}
