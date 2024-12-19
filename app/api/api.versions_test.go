package api_versions_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	apiVersions "github.com/codecrafters-io/kafka-starter-go/api"
)

func TestApiVersionsHandler(t *testing.T) {
	assert.Equal(t, 1, 1)

	request := apiVersions.ApiVersionRequest{
		RequestApiKey:     []byte{0x00, 0x01},
		RequestApiVersion: []byte{0x00, 0x02},
		CorellationId:     []byte{0x12, 0xeb, 0xd1, 0x3c},
		Body:              []byte{0x01, 0x02},
	}

	handler := apiVersions.ApiHandlerApiVersions{}

	resp, err := handler.ProcessRequest(&request)

	assert.Equal(t, resp, nil)
	assert.Nil(t, err)
}
