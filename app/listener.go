package main

import (
	"encoding/binary"
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

		context := &Context{conn: conn}
		listener.registerChannel <- context

		go listener.handleRequest(context)
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

func (listener InnfisListener) handleRequest(context *Context) {
	defer func() {
		listener.unregisterChannel <- context
		context.conn.Close()
	}()

	for {
		message, err := listener.byteToMessage(context)
		if err != nil {
			fmt.Println("handleRequest error: ", err)
			break
		}

		// TODO: what is the code for each api?

		outBuf := ToVersionsResponse(message)
		_, writeErr := context.conn.Write(outBuf)
		if writeErr != nil {
			fmt.Println("handleRequest error: ", writeErr.Error())
			break
		}
	}
}

func (listern InnfisListener) byteToMessage(context *Context) (*Message, error) {
	requestLengthBuffer := make([]byte, 4)
	_, readErr := context.conn.Read(requestLengthBuffer)
	if readErr != nil {
		return nil, readErr
	}

	payloadLen := binary.BigEndian.Uint32(requestLengthBuffer)
	payloadBuffer := make([]byte, payloadLen)
	_, readErr = context.conn.Read(payloadBuffer)
	if readErr != nil {
		return nil, readErr
	}

	message := ToMessage(payloadBuffer)
	return &message, nil
}
