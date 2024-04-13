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
	ctxTimeout, cancel := context.WithTimeout(context.Background(), _defaultTimeout)
	defer cancel()

	resCh := make(chan []*entity.Transaction, 1)

	errCh := make(chan error, 1)

	go w.getTransactionsByBlockNumber(blockNumber, resCh, errCh)

	select {
	case <-ctxTimeout.Done():
		return nil, entity.ErrProcessTimeout
	case <-ctx.Done():
		return nil, entity.ErrInternalServer
	case err := <-errCh:
		return nil, err
	case res := <-resCh:
		return res, nil
	}
}

// Getting current block number from getblock.io.
func (w *StatsOfChangingWebAPI) GetCurrentBlockNumber(ctx context.Context) (*big.Int, error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), _defaultTimeout)
	defer cancel()

	resCh := make(chan *big.Int, 1)

	errCh := make(chan error, 1)

	go w.getCurrentBlockNumber(resCh, errCh)

	select {
	case <-ctxTimeout.Done():
		return nil, entity.ErrProcessTimeout
	case <-ctx.Done():
		return nil, entity.ErrInternalServer
	case err := <-errCh:
		return nil, err
	case res := <-resCh:
		return res, nil
	}
}
