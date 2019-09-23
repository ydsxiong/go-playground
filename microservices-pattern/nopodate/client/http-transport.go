package client

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

	"github.com/ydsxiong/golang/microservices-pattern/nopodate"

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
)

// create a transport to each endpoint to for the client to call exposed, correponding service
func NewHTTPClientTransport(instance string, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) (nopodate.DateService, error) {
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
	// middlewares to demonstrate how to specialize per-endpoint.
	var getEndpoint endpoint.Endpoint
	{
		getEndpoint = httptransport.NewClient(
			http.MethodGet,
			copyURL(u, "/get"),
			prepareApiRequest,
			processGetResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)))...,
		).Endpoint()
		getEndpoint = opentracing.TraceClient(otTracer, "Get")(getEndpoint)
		getEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Get")(getEndpoint)
		getEndpoint = limiter(getEndpoint)
		getEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Get",
			Timeout: 10 * time.Second,
		}))(getEndpoint)
	}

	var statusEndpoint endpoint.Endpoint
	{
		statusEndpoint = httptransport.NewClient(
			http.MethodGet,
			copyURL(u, "/status"),
			prepareApiRequest,
			processStatusResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)))...,
		).Endpoint()
		statusEndpoint = opentracing.TraceClient(otTracer, "Status")(statusEndpoint)
		statusEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Status")(statusEndpoint)
		statusEndpoint = limiter(statusEndpoint)
		statusEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Status",
			Timeout: 10 * time.Second,
		}))(statusEndpoint)
	}

	var validateEndpoint endpoint.Endpoint
	{
		validateEndpoint = httptransport.NewClient(
			http.MethodPost,
			copyURL(u, "/validate"),
			prepareApiRequest,
			processValidateResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)))...,
		).Endpoint()
		validateEndpoint = opentracing.TraceClient(otTracer, "Validate")(validateEndpoint)
		validateEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Validate")(validateEndpoint)
		validateEndpoint = limiter(validateEndpoint)
		validateEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Validate",
			Timeout: 30 * time.Second,
		}))(validateEndpoint)
	}

	return Endpoints{getEndpoint, statusEndpoint, validateEndpoint}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base // dont want to alter the base in any way
	next.Path = path
	return &next
}

func errorDecoder(r *http.Response) error {
	var w nopodate.ErrorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

func prepareApiRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func processGetResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp nopodate.GetResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func processStatusResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp nopodate.StatusResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func processValidateResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp nopodate.ValidateResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
