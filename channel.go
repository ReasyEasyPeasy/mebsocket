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
	for {
		t, m, err := c.conn.ReadMessage()
		if err != nil {
			closeConnection(c)
			return
		}
		fmt.Printf("Messagetype: %v Message: %v\n", t, string(m))
	}
}
func (c channel) write(m string) {
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(m)); err != nil {
		fmt.Println(err)
		closeConnection(c)
	}
}
