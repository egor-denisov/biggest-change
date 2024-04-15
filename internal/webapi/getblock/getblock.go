package webapi

import (
	"context"
	"math/big"
	"net/http"
	"time"

	"github.com/egor-denisov/biggest-change/internal/entity"
	"github.com/egor-denisov/biggest-change/pkg/limiter"
)

var _defaultTimeWindow = 1000 * time.Millisecond
var _defaultTimeout = 15 * time.Second

type StatsOfChangingWebAPI struct {
	URL     string
	Client  *http.Client
	Limiter *limiter.RPSLimiter
}

func New(URL string, rps int) *StatsOfChangingWebAPI {
	if !isValidUrl(URL) {
		panic("invalid getblock.io url: " + URL)
	}

	return &StatsOfChangingWebAPI{
		URL:     URL,
		Client:  &http.Client{},
		Limiter: limiter.NewRPSLimiter(rps, _defaultTimeWindow),
	}
}

// Getting transactions by block number from getblock.io.
func (w *StatsOfChangingWebAPI) GetTransactionsByBlockNumber(
	ctx context.Context,
	blockNumber *big.Int,
) ([]*entity.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	return w.getTransactionsByBlockNumber(ctx, blockNumber)
}

// Getting current block number from getblock.io.
func (w *StatsOfChangingWebAPI) GetCurrentBlockNumber(ctx context.Context) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	return w.getCurrentBlockNumber(ctx)
}
