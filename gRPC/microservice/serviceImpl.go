package microservice

import (
	"github.com/ydsxiong/go-playground/gRPC/microservice/endpoint"
	"github.com/ydsxiong/go-playground/gRPC/microservice/server"
)

type ServiceMgr struct {
	// A map of all the servers registered with the service
	servers map[server.ServerType]server.ServerInterface
}

func (ss *ServiceMgr) RegisterServer(serverType server.ServerType, server server.ServerInterface) (err error) {
	// Store the server in a map
	ss.servers[serverType] = server
	return
}
func (ss *ServiceMgr) Start() (err error) {
	// Do all your service initialization here.
	err = ss.init()
	if err != nil {
		return
	}
	// Invoke the start method on each of the servers
	for t, s := range ss.servers {
		err = s.Start()
		if err != nil {
			return err
		}
	}
	return
}
func (ss *ServiceMgr) Stop() (err error) {
	// Issue a graceful shutdown call to each of the servers
	for t, s := range ss.servers {
		err = s.Stop()
		if err != nil {
			return err
		}
	}
	return
}

// Do all other initialization work, specific to the microservice here.
func (ss *ServiceMgr) init() (err error) {
	sv := ss.servers[server.ServerType("kafka")]
	// Register a namespace
	sv.RegisterNamespace("greeter")

	// Register your endpoints against the namespace with the server
	sv.RegisterService("greeter", server.Service("hello"), endpoint.NewHelloEndpoint())
	return
}
