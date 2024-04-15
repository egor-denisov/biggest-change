package webapi

import (
	"time"
)

type Option func(*StatsOfChangingWebAPI)

func RequestCountRPS(requestCountRPS int) Option {
	return func(s *StatsOfChangingWebAPI) {
		s.requestCountRPS = requestCountRPS
	}
}

func TimeWindowRPS(timeWindowRPS time.Duration) Option {
	return func(s *StatsOfChangingWebAPI) {
		s.timeWindowRPS = timeWindowRPS
	}
}

func Timeout(timeout time.Duration) Option {
	return func(s *StatsOfChangingWebAPI) {
		s.timeout = timeout
	}
}

func MaxRetries(maxRetries int) Option {
	return func(s *StatsOfChangingWebAPI) {
		s.maxRetries = maxRetries
	}
}

func TimeBetweenRetries(timeBetweenRetries time.Duration) Option {
	return func(s *StatsOfChangingWebAPI) {
		s.timeBetweenRetries = timeBetweenRetries
	}
}
