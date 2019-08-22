package handler

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
	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/stringservice/service"
)

type Client struct {
	ctx               context.Context
	UppercaseEndpoint endpoint.Endpoint
	CountEndpoint     endpoint.Endpoint
}

func (c Client) Uppercase(s string) (string, error) {
	response, err := c.UppercaseEndpoint(c.ctx, uppercaseRequest{Input: s})
	if err != nil {
		return "", err
	}

	resp := response.(uppercaseResponse)
	if resp.Err != "" {
		return resp.Output, errors.New(resp.Err)
	}
	return resp.Output, nil
}

func (c Client) Count(s string) int {
	response, _ := c.CountEndpoint(c.ctx, countRequest{Input: s})
	resp := response.(countResponse)
	return resp.Output
}

func NewServiceClient(instance string, ctx context.Context, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) (service.StringService, error) {
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
	var uppercaseEndpoint endpoint.Endpoint
	{
		uppercaseEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/uppercase"),
			prepareApiRequest,
			processUppercaseResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)))...,
		).Endpoint()
		uppercaseEndpoint = opentracing.TraceClient(otTracer, "Uppercase")(uppercaseEndpoint)
		uppercaseEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Uppercase")(uppercaseEndpoint)
		uppercaseEndpoint = limiter(uppercaseEndpoint)
		uppercaseEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Uppercase",
			Timeout: 20 * time.Second,
		}))(uppercaseEndpoint)
	}

	var countEndpoint endpoint.Endpoint
	{
		countEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/count"),
			prepareApiRequest,
			processCountResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)))...,
		).Endpoint()
		countEndpoint = opentracing.TraceClient(otTracer, "Count")(countEndpoint)
		countEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Count")(countEndpoint)
		countEndpoint = limiter(countEndpoint)
		countEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Count",
			Timeout: 10 * time.Second,
		}))(countEndpoint)
	}

	// Returning the Client implementing both of the string service methods.
	return Client{ctx, uppercaseEndpoint, countEndpoint}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base // dont want to alter the base in any way
	next.Path = path
	return &next
}

func prepareApiRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func processUppercaseResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp uppercaseResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func processCountResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp countResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
