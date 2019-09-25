package microservice

import "github.com/ydsxiong/go-playground/gRPC/microservice/server"

type Microservice interface {
	// RegisterServer registers a server with a ServerType identifier. The server must implement the ServerInterface
	RegisterServer(serverType server.ServerType, server server.ServerInterface) (err error)
	// Start is used to start the microservice. All servers registered with the microservice are started when this
	// method is invoked.
	Start() (err error)
	// Stop is used to gracefully stop the microservice. All servers registered with the microservice are gracefully stopped
	// when this method is invoked.
	Stop() (err error)
}
