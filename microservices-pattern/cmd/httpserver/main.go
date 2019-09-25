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

/**
var (
      httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
   )
   flag.Parse()

   var logger log.Logger
   {
      logger = log.NewLogfmtLogger(os.Stderr)
      logger = log.NewSyncLogger(logger)
      logger = level.NewFilter(logger, level.AllowDebug())
      logger = log.With(logger,
         "svc", "order",
         "ts", log.DefaultTimestampUTC,
         "caller", log.DefaultCaller,
      )
   }

   level.Info(logger).Log("msg", "service started")
   defer level.Info(logger).Log("msg", "service ended")

   var db *sql.DB
   {
      var err error
      // Connect to the "ordersdb" database
      db, err = sql.Open("postgres",
         "postgresql://shijuvar@localhost:26257/ordersdb?sslmode=disable")
      if err != nil {
         level.Error(logger).Log("exit", err)
         os.Exit(-1)
      }
   }

   // Create Order Service
   var svc order.Service
   {
      repository, err := cockroachdb.New(db, logger)
      if err != nil {
         level.Error(logger).Log("exit", err)
         os.Exit(-1)
      }
      svc = ordersvc.NewService(repository, logger)
   }

   var h http.Handler
   {
      endpoints := transport.MakeEndpoints(svc)
      h = httptransport.NewService(endpoints, logger)
   }

   errs := make(chan error)
   go func() {
      c := make(chan os.Signal)
      signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
      errs <- fmt.Errorf("%s", <-c)
   }()

   go func() {
      level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
      server := &http.Server{
         Addr:    *httpAddr,
         Handler: h,
      }
      errs <- server.ListenAndServe()
   }()

   level.Error(logger).Log("exit", <-errs)
}
*/
