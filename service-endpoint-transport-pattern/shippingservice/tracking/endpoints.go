package tracking

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type trackCargoRequest struct {
	ID string
}

type trackCargoResponse struct {
	Cargo *CargoView `json:"cargo,omitempty"`
	Err   error      `json:"error,omitempty"`
}

func (r trackCargoResponse) error() error { return r.Err }

func makeTrackCargoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return processTrackingRequest(s, request)
	}
}

func processTrackingRequest(s Service, request interface{}) (interface{}, error) {
	req := request.(trackCargoRequest)
	c, err := s.Track(req.ID)
	return trackCargoResponse{Cargo: &c, Err: err}, nil
}
