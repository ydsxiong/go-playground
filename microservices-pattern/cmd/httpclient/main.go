package main

import (
	"context"
	"flag"
	"log"
	"os"

	gokitlog "github.com/go-kit/kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/ydsxiong/golang/microservices-pattern/nopodate/client"
)

/**
create a http client to make service request call to http server
*/
func main() {
	var (
		serverHost = flag.String("serverHost", "localhost:8080", "remote server host address")
	)
	flag.Parse()

	// Logging domain.
	var logger gokitlog.Logger
	{
		logger = gokitlog.NewLogfmtLogger(os.Stderr)
		logger = gokitlog.With(logger, "ts", gokitlog.DefaultTimestampUTC)
		logger = gokitlog.With(logger, "caller", gokitlog.DefaultCaller)
	}

	// Transport domain.
	tracer := stdopentracing.GlobalTracer() // by default, it's no-op
	zipkinTracer, _ := stdzipkin.NewTracer(nil, stdzipkin.WithNoopTracer(true))
	ctx := context.Background()

	httpClient, err := client.NewHTTPClientTransport(*serverHost, tracer, zipkinTracer, logger)
	if err != nil {
		log.Fatalf("Unable to create a http client, %v", err)
	}

	statusresult, err := httpClient.Status(ctx)
	if err != nil {
		log.Println("Error received from Status call, %v", err)
	} else {
		log.Printf("received from http client Status call: %s", statusresult)
	}

	getresult, err := httpClient.Get(ctx)
	if err != nil {
		log.Println("Error received from Get call, %v", err)
	} else {
		log.Printf("received from http client Get call: %s", getresult)
	}

	validateresult, err := httpClient.Validate(ctx, "26/09/2019")
	if err != nil {
		log.Printf("Error received from Validate call, %v", err)
	} else {
		log.Printf("received from http client Validate call: %t", validateresult)
	}

}
