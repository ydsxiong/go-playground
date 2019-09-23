package nopodate

import (
	"context"
	"errors"
)

/**

At Google, we developed a context package that makes it easy to pass request-scoped values,
cancelation signals, and deadlines across API boundaries to all the goroutines involved in handling a request

Basically, this is needed because our microservice should be made from the beginning to
handle concurrent requests and a context for every request is mandatory.

*/

// Service provides some "date capabilities" to your application
type DateService interface {
	Status(ctx context.Context) (string, error)
	Get(ctx context.Context) (string, error)
	Validate(ctx context.Context, date string) (bool, error)
}

/**
API request response for exposing those operations defined within this service
*/
type GetRequest struct{}

type GetResponse struct {
	Date string `json:"date"`
	Err  string `json:"err,omitempty"`
}

type ValidateRequest struct {
	Date string `json:"date"`
}

type ValidateResponse struct {
	Valid bool   `json:"valid"`
	Err   string `json:"err,omitempty"`
}

type StatusRequest struct{}

type StatusResponse struct {
	Status string `json:"status"`
}

type ErrorWrapper struct {
	Error string `json:"error"`
}

var (
	ErrIncorrectRequestData = errors.New("Wrong input data for api requests!")
)
