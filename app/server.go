package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

var (
	contexts   = make(map[*Context]bool)
	register   = make(chan *Context)
	unregister = make(chan *Context)
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}

	go handleMessage()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}

		context := &Context{conn: conn}
		register <- context
		go handleContext(context)
	}
}

func handleMessage() {
	for {
		select {
		case context := <-register:
			fmt.Println("received register by channel")
			contexts[context] = true
		case context := <-unregister:
			fmt.Println("received unregister by channel")
			delete(contexts, context)
		}
	}
}

func handleContext(context *Context) {
	defer func() {
		unregister <- context
		context.conn.Close()
	}()

	for {
		lenBuf := make([]byte, 4)
		_, readErr := context.conn.Read(lenBuf)
		if readErr != nil {
			fmt.Println(readErr)
			break
		}

		payloadLen := binary.BigEndian.Uint32(lenBuf)

		payloadBuf := make([]byte, payloadLen)
		_, readErr = context.conn.Read(payloadBuf)
		if readErr != nil {
			fmt.Println(readErr)
			break
		}

		message := ToMessage(payloadBuf)
		outBuf := ToOutputBuffer(message)

		_, writeErr := context.conn.Write(outBuf)
		if writeErr != nil {
			fmt.Println("write error: ", writeErr.Error())
		}
	}
}

func ToOutputBuffer(message Message) []byte {
	requestApiVersion := binary.BigEndian.Uint16(message.RequestApiVersion)
	if requestApiVersion > 4 {
		return ToErrorBuffer(message.CorellationId)
	}

	return ToOkBuffer(message.CorellationId)
}
