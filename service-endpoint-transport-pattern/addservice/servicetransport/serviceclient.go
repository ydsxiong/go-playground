package servicetransport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/time/rate"

	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	_ "github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/addservice/service"
	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/addservice/serviceendpoint"
)

// NewServiceClient returns an service backed by an HTTP server living at the
// remote instance. We expect instance to come from a service discovery system,
// so likely of the form "host:port". We bake-in certain middlewares,
// implementing the client library pattern.
func NewServiceClient(instance string, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) (service.AddService, error) {
	// Quickly sanitize the instance string.
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	// Zipkin HTTP Client Trace can either be instantiated per endpoint with a
	// provided operation name or a global tracing client can be instantiated
	// without an operation name and fed to each Go kit endpoint as ClientOption.
	// In the latter case, the operation name will be the endpoint's http method.
	zipkinClient := zipkin.HTTPClientTrace(zipkinTracer)

	// global client middlewares
	options := []httptransport.ClientOption{
		zipkinClient,
	}

	// Each individual endpoint is an http/transport.Client (which implements
	// endpoint.Endpoint) that gets wrapped with various middlewares. If you
	// made your own client library, you'd do this work there, so your server
	// could rely on a consistent set of client behavior.
	var sumEndpoint endpoint.Endpoint
	{
		sumEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/sum"),
			prepareApiRequest,
			processSumResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)))...,
		).Endpoint()
		sumEndpoint = opentracing.TraceClient(otTracer, "Sum")(sumEndpoint)
		sumEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Sum")(sumEndpoint)
		sumEndpoint = limiter(sumEndpoint)
		sumEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Sum",
			Timeout: 30 * time.Second,
		}))(sumEndpoint)
	}

	// The Concat endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var concatEndpoint endpoint.Endpoint
	{
		concatEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/concat"),
			prepareApiRequest,
			processConcatResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)))...,
		).Endpoint()
		concatEndpoint = opentracing.TraceClient(otTracer, "Concat")(concatEndpoint)
		concatEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Concat")(concatEndpoint)
		concatEndpoint = limiter(concatEndpoint)
		concatEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Concat",
			Timeout: 10 * time.Second,
		}))(concatEndpoint)
	}

	// Returning the endpoint.Set implementing the Service methods.
	return serviceendpoint.Set{sumEndpoint, concatEndpoint}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base // dont want to alter the base in any way
	next.Path = path
	return &next
}

// encodeHTTPGenericRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func prepareApiRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// decodeHTTPSumResponse is a transport/http.DecodeResponseFunc that decodes a
// JSON-encoded sum response from the HTTP response body. If the response has a
// non-200 status code, we will interpret that as an error and attempt to decode
// the specific error message from the response body. Primarily useful in a
// client.
func processSumResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp serviceendpoint.SumResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// decodeHTTPConcatResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded concat response from the HTTP response body. If the response
// has a non-200 status code, we will interpret that as an error and attempt to
// decode the specific error message from the response body. Primarily useful in
// a client.
func processConcatResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp serviceendpoint.ConcatResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
