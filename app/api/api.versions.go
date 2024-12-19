package api_versions

import (
	"encoding/binary"
	"fmt"
)

type ApiVersionsRequest struct {
	RequestApiKey     []byte
	RequestApiVersion []byte
	CorellationId     []byte
	Body              []byte
}

type ApiVersionSpec struct {
	ApiKey        uint16
	ApiVersionMin uint16
	ApiVersionMax uint16
}

type ApiVersionsResponse struct {
	Versions     []ApiVersionSpec
	ThrottleTime uint32
}

type ApiHandlerApiVersions struct {
}

func (handler ApiHandlerApiVersions) ProcessRequst(request *ApiVersionsRequest) (*ApiVersionsResponse, error) {
	// TODO
	fmt.Printf("ApiHandlerApiVersion.RequestApiKey] RequestApiKey: %d\n", binary.BigEndian.Uint16(request.RequestApiKey))

	return nil, fmt.Errorf("not implemented")
}
