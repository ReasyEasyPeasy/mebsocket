package mebsocket

import (
	"testing"
	"time"
)

type test struct {
	Topic
	message string
}

func TestSubscribe(t *testing.T) {
	c := subscriber{
		topic: test{
			message: "test",
		},
		conn:  nil,
		timer: &time.Timer{},
	}
	c1 := subscriber{
		topic: test {
			message: "töst",
		},
		conn:  nil,
		timer: &time.Timer{},
	}
	newsDistributerI.subscribe(c)
	newsDistributerI.subscribe(c1)
	time.Sleep(time.Second * 1)
	if len(newsDistributerI.subscriber) != 2 {
		t.Errorf("list shoud be lenght of 2 but was %v", len(newsDistributerI.subscriber))
	}
}

func TestUnsubscribe(t *testing.T) {
	c := subscriber{
		topic: test {
			message: "töst",
		},
		conn:  nil,
		timer: &time.Timer{},
	}
	c1 := subscriber{
		topic: test {
			message: "töst",
		},
		conn:  nil,
		timer: &time.Timer{},
	}
	newsDistributerI.subscribe(c)
	newsDistributerI.subscribe(c1)
	newsDistributerI.unsubscribe(c)
	time.Sleep(time.Second * 1)
	if len(newsDistributerI.subscriber) != 1 {

		t.Errorf("list shoud be lenght of 1 but was %v", len(newsDistributerI.subscriber))
	}
}
