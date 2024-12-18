package test

import (
	"bytes"
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

type ApiVersionSpec struct {
	ApiKey        uint16
	ApiVersionMin uint16
	ApiVersionMax uint16
}

func StructToBytes(obj interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func TestStructToByteArray(t *testing.T) {
	versionSpec := ApiVersionSpec{
		ApiKey:        18,
		ApiVersionMin: 0,
		ApiVersionMax: 4,
	}

	expectedResult := []byte{0x00, 0x12, 0x00, 0x00, 0x00, 0x04}

	byteArray, _ := StructToBytes(versionSpec)

	// not working as expected
	assert.Equal(t, len(expectedResult), len(byteArray))

	for i := 0; i < len(expectedResult); i++ {
		assert.Equal(t, expectedResult[i], byteArray[i])
	}
}
