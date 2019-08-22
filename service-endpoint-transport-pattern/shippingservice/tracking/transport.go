package tracking

import (
	"context"
	"errors"
	"net/http"

	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/shippingservice/servicecommons"
	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/shippingservice/transportcommons"

	"github.com/gorilla/mux"

	kitlog "github.com/go-kit/kit/log"
	kittransport "github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
)

// MakeHandler returns a handler for the tracking service.
func MakeHandler(ts Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(transportcommons.ProcessResponseError),
	}

	trackCargoHandler := kithttp.NewServer(
		makeTrackCargoEndpoint(ts),
		handleTrackCargoRequest,
		transportcommons.HandleResponse,
		opts...,
	)

	r.Handle(servicecommons.TrackCargoPath.String(), trackCargoHandler).Methods("GET")

	return r
}

func handleTrackCargoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.New("bad route")
	}
	return trackCargoRequest{ID: id}, nil
}
