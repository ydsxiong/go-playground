package client

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/ydsxiong/golang/microservices-pattern/nopodate"
)

// aggregating all endpoints for those api services being exposed
type Endpoints struct {
	GetEndpoint      endpoint.Endpoint
	StatusEndpoint   endpoint.Endpoint
	ValidateEndpoint endpoint.Endpoint
}

func (e Endpoints) Get(ctx context.Context) (string, error) {
	req := nopodate.GetRequest{}
	resp, err := e.GetEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	getResp := resp.(nopodate.GetResponse)
	if getResp.Err != "" {
		return "", errors.New(getResp.Err)
	}
	return getResp.Date, nil
}

func (e Endpoints) Status(ctx context.Context) (string, error) {
	req := nopodate.StatusRequest{}
	resp, err := e.StatusEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	statusResp := resp.(nopodate.StatusResponse)
	return statusResp.Status, nil
}

func (e Endpoints) Validate(ctx context.Context, date string) (bool, error) {
	req := nopodate.ValidateRequest{Date: date}
	resp, err := e.ValidateEndpoint(ctx, req)
	if err != nil {
		return false, err
	}
	validateResp := resp.(nopodate.ValidateResponse)
	if validateResp.Err != "" {
		return false, errors.New(validateResp.Err)
	}
	return validateResp.Valid, nil
}
