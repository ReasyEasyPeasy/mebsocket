package mebsocket

import (
	"fmt"
	"github.com/gorilla/websocket"
)

func removeChannel(c subscriber) {
	if err := c.conn.Close(); err != nil {
		fmt.Println(err)
	}
	newsDistributerI.unsubscribe(c)
}

func closeConnection(conn *websocket.Conn) {
	if err := conn.Close(); err != nil {
		fmt.Println(err)
	}
}
