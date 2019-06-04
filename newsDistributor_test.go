package mebsocket

import (
	"testing"
	"time"
)

type test struct {
	Topic
	teamname string
}

func TestSubscribe(t *testing.T) {
	c := subscriber{
		topic: test{
			teamname: "test",
		},
		conn:  nil,
		timer: &time.Timer{},
	}
	c1 := subscriber{
		topic: test{
			teamname: "töst",
		},
		conn:  nil,
		timer: &time.Timer{},
	}
	distributer.subscribe(c)
	distributer.subscribe(c1)
	time.Sleep(time.Second * 1)
	if len(distributer.subscribers) != 2 {
		t.Errorf("list shoud be lenght of 2 but was %v", len(distributer.subscribers))
	}
}

func TestUnsubscribe(t *testing.T) {
	c := subscriber{
		topic: test{
			teamname: "töst",
		},
		conn:  nil,
		timer: &time.Timer{},
	}
	c1 := subscriber{
		topic: test{
			teamname: "töst",
		},
		conn:  nil,
		timer: &time.Timer{},
	}
	distributer.subscribe(c)
	distributer.subscribe(c1)
	distributer.unsubscribe(c)
	time.Sleep(time.Second * 1)
	if len(distributer.subscribers) != 1 {

		t.Errorf("list shoud be lenght of 1 but was %v", len(distributer.subscribers))
	}
}
