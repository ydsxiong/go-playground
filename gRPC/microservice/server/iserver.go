package server

import "github.com/go-kit/kit/endpoint"

/**

https://medium.com/codezillas/building-a-microservice-framework-in-golang-dd3c9530dff9

Server interface
Each server would need to implement the server interface which is composed of the following :

AStart() method which would start the would start the server as a go routine.

A Stop() method which would gracefully stop the server by sending a message on its quit channel returned at the time of starting

A RegisterNamespace() and RegisterService() to register a namespace which allows for using the same channel/topic
for sending messages destined for different endpoints identified by their service type string.


*/

// ServerType is a typedef for server identifier
type ServerType string

// Service is a typedef for service identifier
type Service string

// ServiceEndpointMap is a typedef for a map of Service endpoints identified by their service type
type ServiceEndpointMap map[Service]endpoint.Endpoint

// Handler is a func signature implemented by the Endpoint Handle() method
type Handler func(in []byte, serviceEndpointMap ServiceEndpointMap) (err error)

// StartStopInterface is composed of the Start() and Stop() methods to be implemented by the server
type StartStopInterface interface {
	// Start is used to start the server
	Start() (err error)
	// Stop is used to gracefully stop the server
	Stop() (err error)
}

// ServerInterface defines the methods that all servers registed with the microservice must implement
type ServerInterface interface {
	StartStopInterface
	// RegisterNamesapce is used to register a namespace (like a kafka channel/topic or grpc namespace)
	// with the server.
	RegisterNamespace(namespace string)
	// RegisterService is used to register a service and its endpoints with the server.
	RegisterService(namespace string, service Service, endpoint endpoint.Endpoint)
}
