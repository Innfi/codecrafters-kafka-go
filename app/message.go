package main

import "encoding/binary"

type Message struct {
	RequestApiKey     []byte
	RequestApiVersion []byte
	CorellationId     []byte
	Body              []byte
}

func ToMessage(buf []byte) Message {
	return Message{
		RequestApiKey:     buf[0:2],
		RequestApiVersion: buf[2:4],
		CorellationId:     buf[4:8],
		Body:              buf[8:],
	}
}

func ToErrorBuffer(corellationId []byte) []byte {
	errOutput := make([]byte, 4)
	copy(errOutput, corellationId)

	errorCode := []byte{0x00, 0x23}
	errOutput = append(errOutput, errorCode...)

	outBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(outBuf, uint32(len(errOutput)))
	outBuf = append(outBuf, errOutput...)

	return outBuf
}
