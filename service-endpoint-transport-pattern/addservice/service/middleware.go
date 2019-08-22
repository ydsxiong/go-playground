package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

type serviceMiddleware func(service AddService) AddService

type loggingMiddlewareService struct {
	logger log.Logger
	svc    AddService
}

func NewLoggingMiddlewareService(logger log.Logger) serviceMiddleware {
	return func(service AddService) AddService {
		return loggingMiddlewareService{logger, service}
	}
}

func (ls loggingMiddlewareService) Sum(ctx context.Context, a, b int) (output int, err error) {
	defer func(begin time.Time) {
		_ = ls.logger.Log(
			"method", "sum",
			"a", a, "b", b,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return ls.svc.Sum(ctx, a, b)
}

func (ls loggingMiddlewareService) Concat(ctx context.Context, a, b string) (output string, err error) {
	defer func(begin time.Time) {
		_ = ls.logger.Log(
			"method", "concate",
			"a", a, "b", b,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return ls.svc.Concat(ctx, a, b)
}

type instrumentingMiddlewareService struct {
	ints  metrics.Counter
	chars metrics.Counter
	svc   AddService
}

func NewInstrumentingMiddlewareService(ints, chars metrics.Counter) serviceMiddleware {
	return func(service AddService) AddService {
		return instrumentingMiddlewareService{ints, chars, service}
	}
}

func (is instrumentingMiddlewareService) Sum(ctx context.Context, a, b int) (output int, err error) {
	output, err = is.svc.Sum(ctx, a, b)
	is.ints.Add(float64(output))
	return
}

func (is instrumentingMiddlewareService) Concat(ctx context.Context, a, b string) (output string, err error) {
	output, err = is.svc.Concat(ctx, a, b)
	is.chars.Add(float64(len(output)))
	return
}
