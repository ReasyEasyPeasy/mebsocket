package mebsocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"strings"
	"sync"
	"time"
)

type subscriber struct {
	conn         *websocket.Conn
	topic        Topic
	timer        *time.Timer
	message      string
	messageMutex sync.RWMutex
	messageCh    chan message
	stopCh       chan bool
}

func (c subscriber) pinger() {
	for {
		time.Sleep(time.Second * 2)
		if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			fmt.Printf("Error beim ping:%v\n", err)
			return
		}
	}
}
func (s subscriber) init() {
	go s.reader()
	go s.connHandler()
	go s.pinger()
}
func (c subscriber) stop() {
	c.stopCh <- true
	closeConnection(c.conn)
	distributer.unsubscribe(c.conn)
}

func (c subscriber) readOnce() (string, error) {
	_, m, err := c.conn.ReadMessage()
	if err != nil {
		closeConnection(c.conn)
		return "", err
	}
	return strings.TrimRight(string(m), "\n"), nil
}
func (c subscriber) connHandler() {
	for {
		t, m, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			c.stop()
			return
		}
		c.messageCh <- message{
			t:       t,
			message: string(m),
		}
	}

}
func (c subscriber) reader() {
	for {

		select {
		case m := <-c.messageCh:
			if m.t == websocket.PongMessage {
				fmt.Println("Pong Message received")
				continue
			}
			if m.t == websocket.PingMessage {
				fmt.Println("Ping Message Send")
				if err := c.conn.WriteMessage(websocket.PongMessage, []byte{}); err != nil {
					fmt.Printf("Pong teamname error:%v\n", err)
					return
				}
			}
			c.messageMutex.Lock()
			c.message = m.message
			c.messageMutex.Unlock()
		case _ = <-c.stopCh:
			return
		}

	}

}

func (c subscriber) getLastMessage() string {
	c.messageMutex.RLock()
	defer c.messageMutex.RUnlock()
	return c.message
}
func (c subscriber) write(m string) {
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(m)); err != nil {
		fmt.Println(err)
		c.stop()
	}
}
