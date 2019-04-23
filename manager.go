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
	s.acceptedConnsChannel = make(chan channel)
	s.newConnsChannel = make(chan *websocket.Conn)
	s.declinedConnsChannel = make(chan *websocket.Conn)
	go s.registerHandler()
	go s.declinedConnectionHandler()
}

type MebServer struct {
	newConnsChannel               chan *websocket.Conn
	acceptedConnsChannel          chan channel
	declinedConnsChannel          chan *websocket.Conn
	needAuthorizationConnsChannel chan *websocket.Conn
	Register                      func(string) (bool, Topic, time.Time)
}

type Topic interface {
	EqualTo(Topic) bool
}

type channel struct {
	conn  *websocket.Conn
	topic Topic
	timer *time.Timer
}

func (c *channel) reader() {
	for {
		t, m, err := c.conn.ReadMessage()
		if err != nil {
			c.conn.Close()
			return
		}
		fmt.Printf("Messagetype: %v Message: %v\n", t, string(m))
	}
}

func (s MebServer) declinedConnectionHandler() {
	for {
		select {
		case conn := <-s.declinedConnsChannel:
			if err := conn.WriteMessage(websocket.TextMessage, []byte("{\"type\":\"UNAUTHORIZED\"}")); err != nil {
				fmt.Println(err)
			}
			if err := conn.Close(); err != nil {
				fmt.Println(err)
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
				}
				if ok, topic, tokenEnd := s.Register(strings.TrimRight(string(m), "\n")); ok {
					timer := time.AfterFunc(tokenEnd.Sub(time.Now()), func() {
						s.needAuthorizationConnsChannel <- conn
					})
					c := channel{
						conn:  conn,
						topic: topic,
						timer: timer,
					}
					go c.reader()
					s.acceptedConnsChannel <- c

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
