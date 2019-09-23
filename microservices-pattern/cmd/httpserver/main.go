package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ydsxiong/golang/microservices-pattern/nopodate/server"
)

/**
Channels are the pipes that connect concurrent goroutines. You can send values into channels from one goroutine
and receive those values into another goroutine.

We then create two goroutines. One to stop the server when we press CTRL+C and one that will actually listen
for incoming requests.


*/
func main() {
	var (
		httpAddr = flag.String("http", ":8080", "http listen address")
	)
	flag.Parse()
	ctx := context.Background()
	// our napodate service
	srv := server.NewDefaultDateService()
	errChan := make(chan error)

	// stop the server when we press CTRL+C
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// mapping endpoints
	endpoints := server.Endpoints{
		GetEndpoint:      server.MakeGetEndpoint(srv),
		StatusEndpoint:   server.MakeStatusEndpoint(srv),
		ValidateEndpoint: server.MakeValidateEndpoint(srv),
	}

	// HTTP transport
	go func() {
		log.Println("napodate is listening on port:", *httpAddr)
		handler := server.NewHTTPServerTransport(ctx, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	log.Fatalln(<-errChan)
}
