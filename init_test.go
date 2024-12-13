package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"encoding/binary"
)

func TestByteArray(t *testing.T) {
	output := []byte{}
	input := []byte{0x00, 0x00, 0x00, 0x23}
	output2 := append(output, input...)

	assert.Equal(t, len(output2), 4)
}

func TestAssembleBytes(t *testing.T) {
	input := []byte{0x00, 0x00, 0x00, 0x23}
	inputLen := len(input)

	output := make([]byte, 4)
	binary.BigEndian.PutUint32(output, uint32(inputLen))

	assert.Equal(t, len(output), 4)

	output = append(output, input...)
	assert.Equal(t, len(output), 8)
}
