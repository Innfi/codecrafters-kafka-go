package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	listener, err := NewListener("0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}

	listener.Run()
}
