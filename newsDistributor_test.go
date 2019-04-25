package mebsocket

import (
	"testing"
	"time"
)

func TestSubscribe(t *testing.T) {
	c := channel{
		topic: "test",
		conn:  nil,
		timer: &time.Timer{},
	}
	c1 := channel{
		topic: "test",
		conn:  nil,
		timer: &time.Timer{},
	}
	newsDistributerI.subscribe(c)
	newsDistributerI.subscribe(c1)
	time.Sleep(time.Second * 1)
	if len := len(newsDistributerI.subscriber["test"]); len != 2 {
		t.Errorf("list shoud be lenght of 2 but was %v", len)
	}
}

func TestUnsubscribe(t *testing.T) {
	c := channel{
		topic: "test",
		conn:  nil,
		timer: &time.Timer{},
	}
	c1 := channel{
		topic: "test",
		conn:  nil,
		timer: &time.Timer{},
	}
	newsDistributerI.subscribe(c)
	newsDistributerI.subscribe(c1)
	newsDistributerI.unsubscribe(c)
	time.Sleep(time.Second * 1)
	if len := len(newsDistributerI.subscriber["test"]); len != 1 {

		t.Errorf("list shoud be lenght of 1 but was %v", len)
	}
}
