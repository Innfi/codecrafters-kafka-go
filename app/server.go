package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
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
	conn, connErr := l.Accept()
	if connErr != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go handleConnection(conn, &wg)

	wg.Wait()
	conn.Close()
	// defer conn.Close()

	// buf := make([]byte, 100)
	// _, readErr := conn.Read(buf)
	// if readErr != nil {
	// 	fmt.Println("read error: ", readErr.Error())
	// 	os.Exit(1)
	// }

	// message := toMessage(buf)
	// outBuf := toOutputBuffer(message)

	// _, writeErr := conn.Write(outBuf)
	// if writeErr != nil {
	// 	fmt.Println("write error: ", writeErr.Error())
	// }
}

func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		lenBuf := make([]byte, 4)
		_, readErr := conn.Read(lenBuf)
		if readErr != nil {
			fmt.Println(readErr)
			break
		}

		payloadLen := binary.BigEndian.Uint32(lenBuf)
		fmt.Printf("readLen: %d\n", payloadLen)

		payloadBuf := make([]byte, payloadLen)
		_, readErr = conn.Read(payloadBuf)
		if readErr != nil {
			fmt.Println(readErr)
			break
		}

		message := toMessage(payloadBuf)
		outBuf := toOutputBuffer(message)

		_, writeErr := conn.Write(outBuf)
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
