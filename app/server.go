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

	requestApiVersion := binary.BigEndian.Uint16(message.RequestApiVersion)

	if requestApiVersion > 4 {
		sendErrorResponse(conn, message)
		os.Exit(0)
	}

	sendOkResponse(conn, message)
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

func toMessage(buf []byte) Message {
	// for index, elem := range buf {
	// 	if index > 13 {
	// 		break
	// 	}
	// 	fmt.Printf("index: %d, elem: %02x\n", index, elem)
	// }

	return Message{
		MessageSize:       buf[0:4],
		RequestApiKey:     buf[4:6],
		RequestApiVersion: buf[6:8],
		CorellationId:     buf[8:12],
	}
}

func sendErrorResponse(conn net.Conn, message Message) {
	errOutput := make([]byte, 4)
	copy(errOutput, message.CorellationId)

	errorCode := []byte{0x00, 0x23}
	errOutput = append(errOutput, errorCode...)

	outBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(outBuf, uint32(len(errOutput)))
	outBuf = append(outBuf, errOutput...)

	_, writeErr := conn.Write(outBuf)
	if writeErr != nil {
		fmt.Println("write error: ", writeErr.Error())
	}
}

func sendOkResponse(conn net.Conn, message Message) {
	okOutput := make([]byte, 4)
	copy(okOutput, message.CorellationId)

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

	_, writeErr := conn.Write(outBuf)
	if writeErr != nil {
		fmt.Println("write error: ", writeErr.Error())
	}
}
