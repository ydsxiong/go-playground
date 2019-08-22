package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/stringservice/service"
)

type uppercaseRequest struct {
	Input string `json:"input"`
}

type uppercaseResponse struct {
	Output string `json:"upper"`
	Err    string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

type countRequest struct {
	Input string `json:"input"`
}

type countResponse struct {
	Output int `json:"length"`
}

// setup the transport operations
func CreateUpdatecaseServiceHandler(svc service.StringService) *httptransport.Server {
	return httptransport.NewServer(
		setupUppercaseServiceEndpoint(svc),
		handleUppercaseRequest,
		handleResponse,
	)
}

func CreateCountServiceHandler(svc service.StringService) *httptransport.Server {
	return httptransport.NewServer(
		setupCountServiceEndpoint(svc),
		handleCountRequest,
		handleResponse,
	)
}

func setupUppercaseServiceEndpoint(svc service.StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		return processUppercaseCall(svc, request)
	}
}

func setupCountServiceEndpoint(svc service.StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		return processCountCall(svc, request)
	}
}

func handleUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func processUppercaseCall(svc service.StringService, request interface{}) (interface{}, error) {
	req := request.(uppercaseRequest)
	result, err := svc.Uppercase(req.Input)
	if err != nil {
		return uppercaseResponse{result, err.Error()}, nil
	}
	return uppercaseResponse{result, ""}, nil
}

func handleCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func processCountCall(svc service.StringService, request interface{}) (interface{}, error) {
	req := request.(countRequest)
	result := svc.Count(req.Input)
	return countResponse{result}, nil
}

func handleResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
