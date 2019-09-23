package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/ydsxiong/golang/microservices-pattern/nopodate"
)

// create a transport to each endpoint to serve its correponding service
func NewHTTPServerTransport(ctx context.Context, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	// chain up some middlewares, if any, like logging, instrumenting, etc. for any of those endpionts...
	//r.Use(...middlewares)

	r.Methods("GET").Path("/status").Handler(httptransport.NewServer(
		endpoints.StatusEndpoint,
		handleStatusRequest,
		handleResponse,
	))

	r.Methods("GET").Path("/get").Handler(httptransport.NewServer(
		endpoints.GetEndpoint,
		handleGetRequest,
		handleResponse,
	))

	r.Methods("POST").Path("/validate").Handler(httptransport.NewServer(
		endpoints.ValidateEndpoint,
		handleValidateRequest,
		handleResponse,
	))

	return r
}

func handleGetRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req nopodate.GetRequest
	return req, nil
}

func handleValidateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req nopodate.ValidateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func handleStatusRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req nopodate.StatusRequest
	return req, nil
}

// Last but not least, we have the encoder for the response output
func handleResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		handleErrorResponse(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func handleErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(nopodate.ErrorWrapper{Error: err.Error()})
}

func err2code(err error) int {
	switch err {
	case nopodate.ErrIncorrectRequestData:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
