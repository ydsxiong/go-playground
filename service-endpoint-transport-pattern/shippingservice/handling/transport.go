package handling

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/shippingservice/servicecommons"
	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/shippingservice/transportcommons"

	"github.com/gorilla/mux"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/shippingservice/cargo"
	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/shippingservice/location"
	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/shippingservice/voyage"
)

var eventTypes = map[string]cargo.HandlingEventType{
	cargo.Receive.String(): cargo.Receive,
	cargo.Load.String():    cargo.Load,
	cargo.Unload.String():  cargo.Unload,
	cargo.Customs.String(): cargo.Customs,
	cargo.Claim.String():   cargo.Claim,
}

// MakeHandler returns a handler for the handling service.
func MakeHandler(hs Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(transportcommons.ProcessResponseError),
	}

	registerIncidentHandler := kithttp.NewServer(
		makeRegisterIncidentEndpoint(hs),
		handleRegisterIncidentRequest,
		transportcommons.HandleResponse,
		opts...,
	)

	r.Handle(servicecommons.RegisterIncidentPath.String(), registerIncidentHandler).Methods("POST")

	return r
}

func handleRegisterIncidentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		CompletionTime time.Time `json:"completion_time"`
		TrackingID     string    `json:"tracking_id"`
		VoyageNumber   string    `json:"voyage"`
		Location       string    `json:"location"`
		EventType      string    `json:"event_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return registerIncidentRequest{
		CompletionTime: body.CompletionTime,
		ID:             cargo.TrackingID(body.TrackingID),
		Voyage:         voyage.Number(body.VoyageNumber),
		Location:       location.UNLocode(body.Location),
		EventType:      eventTypes[body.EventType],
	}, nil
}
