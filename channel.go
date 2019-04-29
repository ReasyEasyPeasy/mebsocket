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

func (c channel) reader() {
	defer removeChannel(c)
	for {
		t, m, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		if t == websocket.PingMessage {
			if err := c.conn.WriteMessage(websocket.PongMessage, m); err != nil {
				fmt.Println(err)
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
