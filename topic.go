package mebsocket

type Topic interface {
	Equal(Topic) bool
}

