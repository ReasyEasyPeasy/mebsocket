package mebsocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
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
	s.newConnsChannel = make(chan *websocket.Conn)
	s.declinedConnsChannel = make(chan *websocket.Conn)
	s.needAuthorizationConnsChannel = make(chan channel)
	go s.registerHandler()
	go s.declinedConnectionHandler()
	go s.needAuthorizationHandler()
}

type MebServer struct {
	newConnsChannel               chan *websocket.Conn
	declinedConnsChannel          chan *websocket.Conn
	needAuthorizationConnsChannel chan channel
	Register                      func(string) (bool, interface{}, time.Time)
}

func (s MebServer) SendMessageForTopic(topic interface{}, m string) {
	newsDistributerI.sendMessageToTopic(topic, m)
}

func (s MebServer) needAuthorizationHandler() {
	for {
		select {
		case c := <-s.needAuthorizationConnsChannel:
			fmt.Println("need authoriation")
			newsDistributerI.unsubscribe(c)
			c.write("{\"type\":\"UNAUTHORIZED\"}")
			removeChannel(c)
		}
	}
}

func (s MebServer) declinedConnectionHandler() {
	for {
		select {
		case conn := <-s.declinedConnsChannel:
			if err := conn.WriteMessage(websocket.TextMessage, []byte("{\"type\":\"UNAUTHORIZED\"}")); err != nil {
				fmt.Println(err)
				conn.Close()
			} else {

				s.newConnsChannel <- conn
			}

		}
	}
}

func (s MebServer) HandleWebsocketRequest(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgradeToWebsocket(w, r)
	if err != nil {
		return err
	}
	s.newConnsChannel <- conn
	return nil
}

func (s MebServer) registerHandler() {
	for {
		select {
		case conn := <-s.newConnsChannel:
			go func() {
				_, m, err := conn.ReadMessage()
				if err != nil {
					fmt.Println(err)
					_ = conn.Close()
					return
				}
				if ok, topic, tokenEnd := s.Register(strings.TrimRight(string(m), "\n")); ok {
					c := channel{
						conn:  conn,
						topic: topic,
						timer: nil,
					}
					timer := time.AfterFunc(tokenEnd.Sub(time.Now()), func() {
						fmt.Println("START authorization")
						s.needAuthorizationConnsChannel <- c
					})

					c.timer = timer
					go c.reader()
					go c.pinger()
					newsDistributerI.subscribe(c)

				} else {
					s.declinedConnsChannel <- conn

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
