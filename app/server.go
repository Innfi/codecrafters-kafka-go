package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

type Context struct {
	conn net.Conn
}

var (
	contexts   = make(map[*Context]bool)
	register   = make(chan *Context)
	unregister = make(chan *Context)
)

type Message struct {
	RequestApiKey     []byte
	RequestApiVersion []byte
	CorellationId     []byte
	Body              []byte
}

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

		message := toMessage(payloadBuf)
		outBuf := toOutputBuffer(message)

		_, writeErr := context.conn.Write(outBuf)
		if writeErr != nil {
			fmt.Println("write error: ", writeErr.Error())
		}
	}
}

func toOutputBuffer(message Message) []byte {
	requestApiVersion := binary.BigEndian.Uint16(message.RequestApiVersion)
	if requestApiVersion > 4 {
		return toErrorBuffer(message.CorellationId)
	}

	return toOkBuffer(message.CorellationId)
}

func toMessage(buf []byte) Message {
	return Message{
		RequestApiKey:     buf[0:2],
		RequestApiVersion: buf[2:4],
		CorellationId:     buf[4:8],
		Body:              buf[8:],
	}
}

func toErrorBuffer(corellationId []byte) []byte {
	errOutput := make([]byte, 4)
	copy(errOutput, corellationId)

	errorCode := []byte{0x00, 0x23}
	errOutput = append(errOutput, errorCode...)

	outBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(outBuf, uint32(len(errOutput)))
	outBuf = append(outBuf, errOutput...)

	return outBuf
}

func toOkBuffer(corellationId []byte) []byte {
	okOutput := make([]byte, 4)
	copy(okOutput, corellationId)

	okOutput = append(okOutput, []byte{0x00, 0x00}...)
	okOutput = append(okOutput, []byte{
		0x04,       // array length
		0x00, 0x01, // api key
		0x00, 0x00, // min supported api version
		0x00, 0x11, // max supported api version
		0x00,       // tag buffer
		0x00, 0x12, // api key
		0x00, 0x00, // min supported api version
		0x00, 0x04, // max supported api version
		0x00,       // tag buffer
		0x00, 0x4b, // api key
		0x00, 0x00, // min supported api version
		0x00, 0x00, // max supported api version
		0x00,                   // tag buffer
		0x00, 0x00, 0x00, 0x00, // throttle time
		0x00, // tag buffer
	}...)

	outBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(outBuf, uint32(len(okOutput)))
	outBuf = append(outBuf, okOutput...)

	return outBuf
}

// func printBuf(buf []byte, readLen int) {
// 	for index, elem := range buf {
// 		if index >= readLen {
// 			break
// 		}
// 		fmt.Printf("%02x ", elem)
// 	}
// 	fmt.Printf("\n")
// }
