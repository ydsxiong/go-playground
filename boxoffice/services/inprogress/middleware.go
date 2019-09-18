package inprogress

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/rcrowley/go-metrics"
	"github.com/ydsxiong/go-playground/boxoffice/model"
)

type serviceMiddleware func(service InProgressGuestService) InProgressGuestService

type loggingMiddlewareService struct {
	logger log.Logger
	svc    InProgressGuestService
}

func NewLoggingMiddlewareService(logger log.Logger) serviceMiddleware {
	return func(service InProgressGuestService) InProgressGuestService {
		return loggingMiddlewareService{logger, service}
	}
}

func (ls loggingMiddlewareService) AddGuestInProgress(guest *model.Guest, reservationTime time.Duration) error {
	defer func(begin time.Time) {
		_ = ls.logger.Log(
			"method", "AddGuestInProgress",
			"guestname", guest.Name,
			"took", time.Since(begin),
		)
	}(time.Now())
	return ls.svc.AddGuestInProgress(guest, reservationTime)
}

func (ls loggingMiddlewareService) RemoveGuestFromInProgress(guest *model.Guest) error {
	defer func(begin time.Time) {
		_ = ls.logger.Log(
			"method", "RemoveGuestFromInProgress",
			"guestname", guest.Name,
			"took", time.Since(begin),
		)
	}(time.Now())
	return ls.svc.RemoveGuestFromInProgress(guest)
}

func (ls loggingMiddlewareService) IsGuestInProcess(guest *model.Guest) (yesorno bool, remaining time.Duration, err error) {
	defer func(begin time.Time) {
		_ = ls.logger.Log(
			"method", "IsGuestInProcess",
			"guestname", guest.Name,
			"running", yesorno,
			"remaining", remaining,
			"took", time.Since(begin),
		)
	}(time.Now())
	return ls.svc.IsGuestInProcess(guest)
}

func (ls loggingMiddlewareService) NumberOfGuestInProcess() (num int, err error) {
	defer func(begin time.Time) {
		_ = ls.logger.Log(
			"method", "NumberOfGuestInProcess",
			"output", num,
			"took", time.Since(begin),
		)
	}(time.Now())
	return ls.svc.NumberOfGuestInProcess()
}

type instrumentingMiddlewareService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	InProgressGuestService
}

func NewInstrumentingMiddlewareService(counter metrics.Counter, latency metrics.Histogram, s InProgressGuestService) InProgressGuestService {
	return &instrumentingMiddlewareService{
		requestCount:           counter,
		requestLatency:         latency,
		InProgressGuestService: s,
	}
}

// TODO: implementing service methods for this instrumenting middleware...
