package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/stringservice/service"

	"golang.org/x/time/rate"

	"github.com/sony/gobreaker"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
)

func CreateProxyingMiddleware(ctx context.Context, instances string, logger log.Logger) service.ServiceMiddleware {
	// If instances is empty, don't proxy.
	if instances == "" {
		logger.Log("proxy_to", "none")
		return func(svc service.StringService) service.StringService { return svc }
	}

	// Set some parameters for our client.
	var (
		qps         = 100                // beyond which we will return an error
		maxAttempts = 3                  // per request, before giving up
		maxTime     = 2500 * time.Second // wallclock time, before giving up
	)

	// Otherwise, construct an serviceendpoint for each instance in the list, and add
	// it to a fixed set of endpoints. In a real service, rather than doing this
	// by hand, you'd probably use package sd's support for your service
	// discovery system.
	var (
		instanceList = strings.Split(instances, ",")
		endpointer   sd.FixedEndpointer
	)
	logger.Log("proxy_to", fmt.Sprint(instanceList))
	for _, instance := range instanceList {
		var e endpoint.Endpoint
		e = setupUppercaseProxyClient(ctx, strings.TrimSpace(instance))
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(e)
		endpointer = append(endpointer, e)
	}

	// Now, build a single, retrying, load-balancing serviceendpoint out of all of
	// those individual endpoints.
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(maxAttempts, maxTime, balancer)

	// And finally, return the ServiceMiddleware, implemented by proxymw.
	return func(svc service.StringService) service.StringService {
		return proxymw{ctx, svc, retry}
	}
}

// proxymw implements StringService, forwarding Uppercase requests to the
// provided serviceendpoint, and serving all other (i.e. Count) requests via the
// next/svc StringService.
type proxymw struct {
	ctx                    context.Context
	svc                    service.StringService // Serve most requests via this service...
	callUppercaseApiSerice endpoint.Endpoint     // ...except Uppercase, which gets served by this serviceendpoint
}

func (mw proxymw) Count(s string) int {
	return mw.svc.Count(s)
}

func (mw proxymw) Uppercase(s string) (string, error) {
	response, err := mw.callUppercaseApiSerice(mw.ctx, uppercaseRequest{Input: s})
	if err != nil {
		return "", err
	}

	resp := response.(uppercaseResponse)
	if resp.Err != "" {
		return resp.Output, errors.New(resp.Err)
	}
	return resp.Output, nil
}

// set up an serviceendpoint for sending off a request and receiving a response on behalf of the proxy client
func setupUppercaseProxyClient(ctx context.Context, instance string) endpoint.Endpoint {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		panic(err)
	}
	if u.Path == "" {
		u.Path = "/uppercase"
	}
	return httptransport.NewClient(
		"GET",
		u,
		prepareApiServiceRequest,
		processApiServiceResponse,
	).Endpoint()
}

// set up the body content for the client request to be sent off to the remote server
func prepareApiServiceRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// extract the data from the response received from the remote server
func processApiServiceResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response uppercaseResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}
