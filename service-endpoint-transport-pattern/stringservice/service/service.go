package service

import (
	"errors"
	"strings"
)

type StringService interface {
	Uppercase(string) (string, error)
	Count(string) int
}

type basicStringService struct{}

func (basicStringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (basicStringService) Count(s string) int {
	return len(s)
}

var ErrEmpty = errors.New("Empty string")

// ServiceMiddleware is a chainable behavior modifier for StringService.
type ServiceMiddleware func(StringService) StringService

func CreateBasicStringService() StringService {
	return basicStringService{}
}