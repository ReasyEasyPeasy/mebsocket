package mebsocket

import (
	"fmt"
)

func closeConnection(c channel) {
	if err := c.conn.Close(); err != nil {
		fmt.Println(err)
	}
	newsDistributerI.unsubscribe(c)
}
