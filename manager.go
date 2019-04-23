package mebsocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *MebServer) Init() {
	fmt.Println("init")
	s.acceptedConnsChannel = make(chan *websocket.Conn)
	s.newConnsChannel = make(chan *websocket.Conn)
	go s.connectionHandler()
}
func NewMebsocket() MebServer {
	s := MebServer{}

	return s
}

type mInterface interface {
	AcceptToken(string) bool
}
type MebServer struct {
	mInterface
	newConnsChannel            chan *websocket.Conn
	acceptedConnsChannel       chan *websocket.Conn
	CheckFirstWebsocketMessage func(string) bool
}

func (s MebServer) AcceptToken(token string) bool {
	fmt.Println("parent")
	return true
}

func (s MebServer) HandleWebsocketRequest(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgradeToWebsocket(w, r)
	if err != nil {
		return err
	}
	s.newConnsChannel <- conn
	return nil
}
func (s MebServer) connectionHandler() {
	for {
		select {
		case conn := <-s.newConnsChannel:
			go func() {
				_, m, err := conn.ReadMessage()
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(string(m))
				if s.AcceptToken(strings.TrimRight(string(m), "\n")) {
					s.acceptedConnsChannel <- conn
				} else {
					fmt.Println("Unauthorized")
					if err := conn.WriteMessage(websocket.TextMessage, []byte("{\"type\":\"UNAUTHORIZED\"}")); err != nil {
						fmt.Println(err)
					}
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
