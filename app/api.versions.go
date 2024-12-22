package main

import "encoding/binary"

func ToOkBuffer(corellationId []byte) []byte {
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
