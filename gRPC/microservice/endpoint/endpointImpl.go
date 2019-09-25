package endpoint

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// HelloEndpoint is the endpoint specific object
type HelloEndpoint struct {
}

// NewHelloEndpoint returns an instance of the HelloEndpoint object
func NewHelloEndpoint() *HelloEndpoint {
	return &HelloEndpoint{}
}

// Handle implements the Endpoint Handle method.
// This is where the business logic of the endpoint is implemented.
func (ep *HelloEndpoint) Handle(ctx context.Context, request []byte) (response []byte, err error) {
	log.Infof("\n=======\nHello %s\n=======\n", request)
	return
}
