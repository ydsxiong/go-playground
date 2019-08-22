package booking

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/shippingservice/cargo"
	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/shippingservice/location"
	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/shippingservice/servicecommons"
	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/shippingservice/transportcommons"

	"github.com/gorilla/mux"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
)

// MakeHandler returns a handler for the booking service.
func MakeHandler(bs Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(transportcommons.ProcessResponseError),
	}

	bookCargoHandler := kithttp.NewServer(
		makeBookCargoEndpoint(bs),
		handleBookCargoRequest,
		transportcommons.HandleResponse,
		opts...,
	)
	loadCargoHandler := kithttp.NewServer(
		makeLoadCargoEndpoint(bs),
		handleLoadCargoRequest,
		transportcommons.HandleResponse,
		opts...,
	)
	requestRoutesHandler := kithttp.NewServer(
		makeRequestRoutesEndpoint(bs),
		handleRequestRoutesRequest,
		transportcommons.HandleResponse,
		opts...,
	)
	assignToRouteHandler := kithttp.NewServer(
		makeAssignToRouteEndpoint(bs),
		handleAssignToRouteRequest,
		transportcommons.HandleResponse,
		opts...,
	)
	changeDestinationHandler := kithttp.NewServer(
		makeChangeDestinationEndpoint(bs),
		handleChangeDestinationRequest,
		transportcommons.HandleResponse,
		opts...,
	)
	listCargosHandler := kithttp.NewServer(
		makeListCargosEndpoint(bs),
		handleListCargosRequest,
		transportcommons.HandleResponse,
		opts...,
	)
	listLocationsHandler := kithttp.NewServer(
		makeListLocationsEndpoint(bs),
		handleListLocationsRequest,
		transportcommons.HandleResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle(servicecommons.CargosPath.String(), bookCargoHandler).Methods("POST")
	r.Handle(servicecommons.CargosPath.String(), listCargosHandler).Methods("GET")
	r.Handle(servicecommons.LoadCargoPath.String(), loadCargoHandler).Methods("GET")
	r.Handle(servicecommons.RequestRoutesPath.String(), requestRoutesHandler).Methods("GET")
	r.Handle(servicecommons.AssignToRoutePath.String(), assignToRouteHandler).Methods("POST")
	r.Handle(servicecommons.ChangeDestinationPath.String(), changeDestinationHandler).Methods("POST")
	r.Handle(servicecommons.ListLocationsPath.String(), listLocationsHandler).Methods("GET")

	return r
}

var errBadRoute = errors.New("bad route")

func handleBookCargoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Origin          string    `json:"origin"`
		Destination     string    `json:"destination"`
		ArrivalDeadline time.Time `json:"arrival_deadline"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return bookCargoRequest{
		Origin:          location.UNLocode(body.Origin),
		Destination:     location.UNLocode(body.Destination),
		ArrivalDeadline: body.ArrivalDeadline,
	}, nil
}

func handleLoadCargoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	return loadCargoRequest{ID: cargo.TrackingID(id)}, nil
}

func handleRequestRoutesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	return requestRoutesRequest{ID: cargo.TrackingID(id)}, nil
}

func handleAssignToRouteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}

	var itinerary cargo.Itinerary
	if err := json.NewDecoder(r.Body).Decode(&itinerary); err != nil {
		return nil, err
	}

	return assignToRouteRequest{
		ID:        cargo.TrackingID(id),
		Itinerary: itinerary,
	}, nil
}

func handleChangeDestinationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}

	var body struct {
		Destination string `json:"destination"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return changeDestinationRequest{
		ID:          cargo.TrackingID(id),
		Destination: location.UNLocode(body.Destination),
	}, nil
}

func handleListCargosRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return listCargosRequest{}, nil
}

func handleListLocationsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return listLocationsRequest{}, nil
}
