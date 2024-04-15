package webapi

import (
	"context"
	"math/big"
	"net/http"
	"time"

	"github.com/egor-denisov/biggest-change/internal/entity"
	"github.com/egor-denisov/biggest-change/pkg/limiter"
)

const (
	_defaultRequestCountRPS    = 60
	_defaultTimeWindowRPS      = 1 * time.Second
	_defaultTimeout            = 15 * time.Second
	_defaultMaxRetries         = 5
	_defaultTimeBetweenRetries = 500 * time.Millisecond
)

type StatsOfChangingWebAPI struct {
	url                string
	client             *http.Client
	limiter            *limiter.RPSLimiter
	requestCountRPS    int
	timeWindowRPS      time.Duration
	timeout            time.Duration
	maxRetries         int
	timeBetweenRetries time.Duration
}

func New(url string, opts ...Option) *StatsOfChangingWebAPI {
	if !isValidUrl(url) {
		panic("invalid getblock.io url: " + url)
	}

	w := &StatsOfChangingWebAPI{
		url:                url,
		client:             &http.Client{},
		requestCountRPS:    _defaultRequestCountRPS,
		timeWindowRPS:      _defaultTimeWindowRPS,
		timeout:            _defaultTimeout,
		maxRetries:         _defaultMaxRetries,
		timeBetweenRetries: _defaultTimeBetweenRetries,
	}

	for _, opt := range opts {
		opt(w)
	}

	w.limiter = limiter.NewRPSLimiter(w.requestCountRPS, w.timeWindowRPS)

	return w
}

// Getting transactions by block number from getblock.io.
func (w *StatsOfChangingWebAPI) GetTransactionsByBlockNumber(
	ctx context.Context,
	blockNumber *big.Int,
) ([]*entity.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	return w.getTransactionsByBlockNumber(ctx, blockNumber)
}

// Getting current block number from getblock.io.
func (w *StatsOfChangingWebAPI) GetCurrentBlockNumber(ctx context.Context) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(ctx, w.timeout)
	defer cancel()

	return w.getCurrentBlockNumber(ctx)
}
