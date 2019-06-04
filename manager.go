package mebsocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *MebServer) Init() {
	s.declinedConnsChannel = make(chan *subscriber)
	s.needAuthorizationConnsChannel = make(chan *subscriber)
	go s.authorizationHandler()

}

type MebServer struct {
	declinedConnsChannel          chan *subscriber
	needAuthorizationConnsChannel chan *subscriber
	Register                      func(string) (bool, Topic, time.Time)
}

func (s MebServer) PublishMessage(topic Topic, m string) {
	distributer.sendMessageToTopic(topic, m)
}

func sendAuthorizationRequest(s *subscriber) {
	s.write("{\"type\":\"UNAUTHORIZED\"}")
}

func (s MebServer) HandleWebsocketRequest(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgradeToWebsocket(w, r)
	if err != nil {
		return err
	}
	s.needAuthorizationConnsChannel <- &subscriber{
		conn:         conn,
		timer:        nil,
		messageMutex: sync.RWMutex{},
		stopCh:       make(chan bool),
		messageCh:    make(chan message),
	}
	return nil
}

func (s MebServer) authorizationHandler() {
	for {
		select {
		case sub := <-s.needAuthorizationConnsChannel:
			go func() {
				message, err := sub.readOnce()
				if err != nil {
					s.declinedConnsChannel <- sub
					return
				}
				if ok, topic, tokenEnd := s.Register(message); ok {

					timer := time.AfterFunc(tokenEnd.Sub(time.Now()), func() {
						sub.stop()
						s.needAuthorizationConnsChannel <- sub
					})

					sub.timer = timer
					sub.topic = topic
					sub.init()
					distributer.subscribe(sub)

				} else {
					sendAuthorizationRequest(sub)
					s.needAuthorizationConnsChannel <- sub

				}
			}()

		}
	}
}

func upgradeToWebsocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
