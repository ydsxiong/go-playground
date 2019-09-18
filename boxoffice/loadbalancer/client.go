package loadbalancer

import "net/http"

type ClientRequest struct {
	Request *http.Request
	Resp    chan http.ResponseWriter
}

func CreateRequest(req *http.Request) ClientRequest {
	return ClientRequest{req, make(chan http.ResponseWriter)}
}
