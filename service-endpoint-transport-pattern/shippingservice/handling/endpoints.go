package handling

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"

	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/shippingservice/cargo"
	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/shippingservice/location"
	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/shippingservice/voyage"
)

type registerIncidentRequest struct {
	ID             cargo.TrackingID
	Location       location.UNLocode
	Voyage         voyage.Number
	EventType      cargo.HandlingEventType
	CompletionTime time.Time
}

type registerIncidentResponse struct {
	Err error `json:"error,omitempty"`
}

func (r registerIncidentResponse) error() error { return r.Err }

func makeRegisterIncidentEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return processRegisterRrequest(svc, request)
	}
}

func processRegisterRrequest(svc Service, request interface{}) (interface{}, error) {
	req := request.(registerIncidentRequest)
	err := svc.RegisterHandlingEvent(req.CompletionTime, req.ID, req.Voyage, req.Location, req.EventType)
	return registerIncidentResponse{Err: err}, nil
}
