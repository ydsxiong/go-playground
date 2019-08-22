package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/stringservice/handler"
	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/stringservice/instrument"
	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/stringservice/logging"
	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/stringservice/service"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
		proxy  = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy uppercase requests")
	)
	flag.Parse()

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)

	svc := service.CreateBasicStringService()
	svc = handler.CreateProxyingMiddleware(context.Background(), *proxy, logger)(svc)
	svc = logging.CreateLoggingMiddleware(logger)(svc)
	svc = instrument.CreateInstrumentingMiddleware()(svc)

	// wire it up all together and fire up the server
	http.Handle("/uppercase", handler.CreateUpdatecaseServiceHandler(svc))
	http.Handle("/count", handler.CreateCountServiceHandler(svc))
	http.Handle("/metrics", promhttp.Handler())
	//logger.Log("msg", "HTTP", "addr", *listen)
	//logger.Log("err", http.ListenAndServe(*listen, nil))
}
