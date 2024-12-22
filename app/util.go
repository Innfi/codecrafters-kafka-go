package main

import "fmt"

func PrintBuf(buf []byte, readLen int) {
	for index, elem := range buf {
		if index >= readLen {
			break
		}
		fmt.Printf("%02x ", elem)
	}
	fmt.Printf("\n")
}
