package main

import "encoding/binary"

func ToDescribeTopicPartitionsResponse(message *Message) []byte {
	okOutput := make([]byte, 4)
	copy(okOutput, message.CorellationId)

	okOutput = append(okOutput, []byte{0x00}...)
	okOutput = append(okOutput, []byte{
		0x00,                   // tag buffer
		0x00, 0x00, 0x00, 0x00, // throttle time
		0x02,                   // aray length
		0x00,                   // error code
		0x04, 0x66, 0x6f, 0x6f, // topic name
		0xe3, 0x92, 0x80, 0x6d, 0xb5, 0x33, 0x48, 0x10, 0xba, 0x03, 0x0b, 0x43, 0xc4, 0x9b, 0x60, 0xfc, // topic id
		0x00, // isInternal
		0x03, // partitions array length
		// TODO: partition array
		0x00, 0x00, 0x0d, 0xf8, // topic authorized operations
		0x00, // tag buffer
		0xff, // next cursor
		0x00, // tag buffer
	}...)

	outBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(outBuf, uint32(len(okOutput)))
	outBuf = append(outBuf, okOutput...)

	return okOutput
}
