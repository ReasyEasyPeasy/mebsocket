package mebsocket

import (
	"fmt"
	"github.com/gorilla/websocket"
)

func closeConnection(conn *websocket.Conn) {
	if err := conn.Close(); err != nil {
		fmt.Println(err)
	}
}
