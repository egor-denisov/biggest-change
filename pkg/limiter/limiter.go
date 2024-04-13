package limiter

import (
	"sync"
	"time"
)

type RPSLimiter struct {
	mu        sync.Mutex
	lastReset time.Time
	window    time.Duration
	cap       int
	count     int
}

func NewRPSLimiter(reqCount int, window time.Duration) *RPSLimiter {
	limiter := &RPSLimiter{
		cap:       reqCount - 1,
		window:    window,
		lastReset: time.Now(),
	}
	go limiter.resetCounter()

	return limiter
}

func (r *RPSLimiter) resetCounter() {
	for {
		time.Sleep(time.Until(r.lastReset.Add(r.window)))
		r.mu.Lock()
		r.count = 0
		r.lastReset = time.Now()
		r.mu.Unlock()
	}
}

func (r *RPSLimiter) WaitForAvailability() {
	for {
		r.mu.Lock()
		if r.count < r.cap {
			r.count++
			r.mu.Unlock()
			time.Sleep(time.Millisecond * 100)
			break
		}
		r.mu.Unlock()
		time.Sleep(time.Millisecond * 100)
	}
}
