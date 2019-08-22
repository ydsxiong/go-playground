package service

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

type AddService interface {
	Sum(ctx context.Context, a, b int) (int, error)
	Concat(ctx context.Context, a, b string) (string, error)
}

func New(logger log.Logger, ints, chars metrics.Counter) AddService {
	var svc AddService
	{
		svc = newBasicService()
		svc = NewLoggingMiddlewareService(logger)(svc)
		svc = NewInstrumentingMiddlewareService(ints, chars)(svc)
	}

	return svc
}

var (
	ErrTwoZeros        = errors.New("Can not sum two zeros!")
	ErrIntOverflow     = errors.New("integer overflow")
	ErrMaxSizeExceeded = errors.New("result exceeds maximum size")
)

func newBasicService() AddService {
	return basicService{}
}

type basicService struct{}

const (
	intMax = 1<<31 - 1
	intMin = -(intMax + 1)
	maxLen = 10
)

func (basicService) Sum(ctx context.Context, a, b int) (int, error) {
	if a == 0 && b == 0 {
		return 0, ErrTwoZeros
	}

	if (a > 0 && b > intMax-a) || (a < 0 && b > intMin-a) {
		return 0, ErrIntOverflow
	}
	return a + b, nil
}

func (basicService) Concat(ctx context.Context, a, b string) (string, error) {
	if len(a) > maxLen || len(b) > maxLen {
		return "", ErrMaxSizeExceeded
	}
	return a + b, nil
}
