package mebsocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

type channel struct {
	conn  *websocket.Conn
	topic interface{}
	timer *time.Timer
}

func (c channel) pinger() {
	for {
		time.Sleep(time.Second * 2)
		if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			fmt.Printf("Error beim ping:%v\n", err)
			return
		}
	}
}
func (c channel) reader() {
	defer removeChannel(c)
	for {
		t, m, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Printf("Read Error: %v \n", err)
			return

		}
		if t == websocket.PongMessage {
			fmt.Println("Pong Message Send")
			continue
		}
		if t == websocket.PingMessage {
			fmt.Println("Ping Message Send")
			if err := c.conn.WriteMessage(websocket.PongMessage, m); err != nil {
				fmt.Printf("Pong message error:%v\n", err)
				return
			}
		}

		fmt.Printf("Messagetype: %v Message: %v\n", t, string(m))
	}
}
func (c channel) write(m string) {
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(m)); err != nil {
		fmt.Println(err)
		removeChannel(c)
	}
}
