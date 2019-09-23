package server

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/ydsxiong/golang/microservices-pattern/nopodate"
)

// aggregating all endpoints for those api services being exposed
type Endpoints struct {
	GetEndpoint      endpoint.Endpoint
	StatusEndpoint   endpoint.Endpoint
	ValidateEndpoint endpoint.Endpoint
}

// MakeGetEndpoint returns the response from our service "get"
func MakeGetEndpoint(srv nopodate.DateService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return processGetRequest(ctx, request, srv)
	}
}

// MakeStatusEndpoint returns the response from our service "status"
func MakeStatusEndpoint(srv nopodate.DateService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return processStatusRequest(ctx, request, srv)
	}
}

// MakeValidateEndpoint returns the response from our service "validate"
func MakeValidateEndpoint(srv nopodate.DateService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return processValidateRequest(ctx, request, srv)
	}
}

func processGetRequest(ctx context.Context, request interface{}, srv nopodate.DateService) (interface{}, error) {
	_ = request.(nopodate.GetRequest) // we really just need the request, we don't use any value from it
	d, err := srv.Get(ctx)
	if err != nil {
		return nopodate.GetResponse{d, err.Error()}, nil
	}
	return nopodate.GetResponse{d, ""}, nil
}

func processStatusRequest(ctx context.Context, request interface{}, srv nopodate.DateService) (interface{}, error) {
	_ = request.(nopodate.StatusRequest) // we really just need the request, we don't use any value from it
	s, err := srv.Status(ctx)
	if err != nil {
		return nopodate.StatusResponse{s}, err
	}
	return nopodate.StatusResponse{s}, nil
}

func processValidateRequest(ctx context.Context, request interface{}, srv nopodate.DateService) (interface{}, error) {
	req := request.(nopodate.ValidateRequest)
	b, err := srv.Validate(ctx, req.Date)
	if err != nil {
		return nopodate.ValidateResponse{b, err.Error()}, nil
	}
	return nopodate.ValidateResponse{b, ""}, nil
}
