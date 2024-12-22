package main

import (
	"fmt"
	"net"
)

type Context struct {
	conn net.Conn
}

type InnfisListener struct {
	handle            net.Listener
	contexts          map[*Context]bool
	registerChannel   chan *Context
	unregisterChannel chan *Context
}

func NewListener(bindAddress string) (*InnfisListener, error) {
	l, err := net.Listen("tcp", bindAddress)
	if err != nil {
		fmt.Println("failed to listen: ", bindAddress)
		return nil, err
	}

	listener := InnfisListener{
		handle:            l,
		contexts:          make(map[*Context]bool),
		registerChannel:   make(chan *Context),
		unregisterChannel: make(chan *Context),
	}

	return &listener, nil
}

func (listener InnfisListener) Run() {
	go listener.handleChannelMessage()

	for {
		conn, err := listener.handle.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		context := Context{conn: conn}
		listener.registerChannel <- &context
	}
}

func (listener InnfisListener) handleChannelMessage() {
	for {
		select {
		case context := <-listener.registerChannel:
			fmt.Println("received context from registerChannel")
			listener.contexts[context] = true
		case context := <-listener.unregisterChannel:
			fmt.Println("received context from unregisterChannel")
			delete(listener.contexts, context)
		}
	}
}
