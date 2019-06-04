package mebsocket

import (
	"fmt"
	"sync"
)

type newsDistributer struct {
	subscriber []subscriber
	mu         sync.Mutex
}

var newsDistributerI newsDistributer

func init() {
	newsDistributerI.mu = sync.Mutex{}
}

func (n *newsDistributer) subscribe(c subscriber) {
	go func() {
		n.mu.Lock()
		defer func() {
			fmt.Printf("Es gibt %v subscriber!", len(n.subscriber))
		}()
		defer n.mu.Unlock()
		n.subscriber = append(n.subscriber, c)
	}()
}
func (n *newsDistributer) sendMessageToTopic(topic Topic, message string) {
	go func() {
		n.mu.Lock()
		defer n.mu.Unlock()
		for _, c := range n.subscriber {
			if c.topic.Equal(topic) {
				c.write(message)
			}

		}
	}()

}
func (n *newsDistributer) unsubscribe(c subscriber) {

	go func() {
		n.mu.Lock()
		defer func() {
			fmt.Printf("Es gibt %v subscriber!", len(n.subscriber))
		}()
		defer n.mu.Unlock()
		for i, c1 := range n.subscriber {
			if c1 == c {
				n.subscriber = append(n.subscriber[:i], n.subscriber[i+1:]...)
				return
			}
		}

	}()
}
