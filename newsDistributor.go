package mebsocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type newsDistributer struct {
	subscribers []*subscriber
	mu          sync.Mutex
}

var distributer newsDistributer

func init() {
	distributer.mu = sync.Mutex{}
}

func (n *newsDistributer) subscribe(c *subscriber) {
	go func() {
		n.mu.Lock()
		defer func() {
			fmt.Printf("Es gibt %v subscribers!\n", len(n.subscribers))
		}()
		defer n.mu.Unlock()
		n.subscribers = append(n.subscribers, c)
	}()
}
func (n *newsDistributer) sendMessageToTopic(topic Topic, message string) {
	go func() {
		n.mu.Lock()
		defer n.mu.Unlock()
		for _, sub := range n.subscribers {
			if sub.topic.Equal(topic) {
				sub.write(message)
			}

		}
	}()

}
func (n *newsDistributer) unsubscribe(c *websocket.Conn) {

	go func() {
		n.mu.Lock()
		defer func() {
			fmt.Printf("Es gibt %v subscribers!\n", len(n.subscribers))
		}()
		defer n.mu.Unlock()
		for i, sub := range n.subscribers {
			if sub.conn == c {
				n.subscribers = append(n.subscribers[:i], n.subscribers[i+1:]...)
				return
			}
		}

	}()
}
