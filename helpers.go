package mebsocket

import (
	"fmt"
)

func removeChannel(c channel) {
	if err := c.conn.Close(); err != nil {
		fmt.Println(err)
	}
	newsDistributerI.unsubscribe(c)
}
