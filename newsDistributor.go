package mebsocket

import (
	"sync"
)

type newsDistributer struct {
	subscriber map[interface{}][]channel
	mu         sync.Mutex
}

var newsDistributerI newsDistributer

func init() {
	newsDistributerI.subscriber = make(map[interface{}][]channel)
	newsDistributerI.mu = sync.Mutex{}
}
func (n *newsDistributer) subscribe(c channel) {
	go func() {
		n.mu.Lock()
		defer n.mu.Unlock()
		n.subscriber[c.topic] = append(n.subscriber[c.topic], c)
	}()
}
func (n *newsDistributer) sendMessageToTopic(topic interface{}, message string) {
	go func() {
		n.mu.Lock()
		defer n.mu.Unlock()
		for _, c := range n.subscriber[topic] {
			c.write(message)
		}
	}()

}
func (n *newsDistributer) unsubscribe(c channel) {
	go func() {
		n.mu.Lock()
		defer n.mu.Unlock()
		for i, c1 := range n.subscriber[c.topic] {
			if c1 == c {
				n.subscriber[c.topic] = append(n.subscriber[c.topic][:i], n.subscriber[c.topic][i+1:]...)
				return
			}
		}

	}()
}
