package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

type Message struct {
	MessageSize       []byte
	RequestApiKey     []byte
	RequestApiVersion []byte
	CorellationId     []byte
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	conn, connErr := l.Accept()
	if connErr != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	buf := make([]byte, 100)
	_, readErr := conn.Read(buf)
	if readErr != nil {
		fmt.Println("read error: ", readErr.Error())
		os.Exit(1)
	}

	message := toMessage(buf)
	// printBuf(message.MessageSize, 4)
	// printBuf(message.RequestApiKey, 2)
	// printBuf(message.RequestApiVersion, 2)
	// printBuf(message.CorellationId, 4)

	requestApiVersion := binary.BigEndian.Uint16(message.RequestApiVersion)
	// fmt.Printf("requestApiVersion: %d\n", requestApiVersion)

	if requestApiVersion > 4 {
		sendErrorResponse(conn, message)
		os.Exit(0)
	}

	sendOkResponse(conn, message)
}

func printBuf(buf []byte, readLen int) {
	for index, elem := range buf {
		if index >= readLen {
			break
		}
		fmt.Printf("%02x ", elem)
	}
	fmt.Printf("\n")
}

func toMessage(buf []byte) Message {
	for index, elem := range buf {
		if index > 13 {
			break
		}
		fmt.Printf("index: %d, elem: %02x\n", index, elem)
	}

	return Message{
		MessageSize:       buf[0:4],
		RequestApiKey:     buf[4:6],
		RequestApiVersion: buf[6:8],
		CorellationId:     buf[8:12],
	}
}

func sendErrorResponse(conn net.Conn, message Message) {
	outBuf := make([]byte, 10)
	copy(outBuf[4:8], message.CorellationId)
	outBuf[8] = 0x00
	outBuf[9] = 0x23

	printBuf(outBuf, len(outBuf))

	_, writeErr := conn.Write(outBuf)
	if writeErr != nil {
		fmt.Println("write error: ", writeErr.Error())
	}
}

func sendOkResponse(conn net.Conn, message Message) {
	outBuf := make([]byte, 8)
	copy(outBuf[4:8], message.CorellationId)

	printBuf(outBuf, len(outBuf))

	_, writeErr := conn.Write(outBuf)
	if writeErr != nil {
		fmt.Println("write error: ", writeErr.Error())
	}
}
