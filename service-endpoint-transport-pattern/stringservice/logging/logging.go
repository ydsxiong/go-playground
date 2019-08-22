package logging

import (
	"time"

	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/stringservice/service"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	Logger  log.Logger
	Service service.StringService
}

func (mw loggingMiddleware) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.Service.Uppercase(s)
	return
}

func (mw loggingMiddleware) Count(s string) (n int) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "count",
			"input", s,
			"n", n,
			"took", time.Since(begin),
		)
	}(time.Now())

	n = mw.Service.Count(s)
	return
}

func CreateLoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(svc service.StringService) service.StringService {
		return loggingMiddleware{logger, svc}
	}
}
