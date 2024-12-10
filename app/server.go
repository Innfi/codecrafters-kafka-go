package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

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

	fmt.Printf("buffer: %s\n", string(buf))

	outBuf := make([]byte, 8)
	binary.BigEndian.PutUint32(outBuf[4:], 0x07)

	_, writeErr := conn.Write(outBuf)
	if writeErr != nil {
		fmt.Println("write error: ", writeErr.Error())
	}
}
