package usecase

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/egor-denisov/biggest-change/internal/entity"
	lru "github.com/hashicorp/golang-lru"
)

const (
	_defaultMaxGoroutines                   = 50
	_defaultAverageAddressCountInBlock      = 200
	_defaultCacheSize                       = 100 // Count of blocks for which the transaction value will be cached
	_defaultCountOfBlocks              uint = 100
)

type StatsOfChangingUseCase struct {
	webAPI                     StatsOfChangingWebAPI
	cache                      *lru.Cache
	cacheSize                  int
	maxGoroutines              int
	averageAddressCountInBlock int
	countOfBlocks              uint
}

func New(w StatsOfChangingWebAPI, opts ...Option) *StatsOfChangingUseCase {
	uc := &StatsOfChangingUseCase{
		webAPI:                     w,
		cacheSize:                  _defaultCacheSize,
		maxGoroutines:              _defaultMaxGoroutines,
		averageAddressCountInBlock: _defaultAverageAddressCountInBlock,
		countOfBlocks:              _defaultCountOfBlocks,
	}

	for _, opt := range opts {
		opt(uc)
	}

	uc.cache, _ = lru.New(uc.cacheSize)

	return uc
}

// Get address with biggest change in last countOfLastBlocks blocks.
func (uc *StatsOfChangingUseCase) GetAddressWithBiggestChange(
	ctx context.Context,
	countOfLastBlocks uint,
) (*entity.BiggestChange, error) {
	if countOfLastBlocks == 0 {
		countOfLastBlocks = uc.countOfBlocks
	}
	// Getting current number of block.
	currentBlock, err := uc.webAPI.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil,
			fmt.Errorf("StatsOfChangingUseCase - GetAddressWithBiggestChange - uc.webAPI.GetCurrentBlockNumber: %w", err)
	}
	// Getting map which store addresses and changes in last countOfLastBlocks blocks.
	addresses, err := uc.getAddressChangeMap(ctx, currentBlock, int(countOfLastBlocks))
	if err != nil {
		return nil,
			fmt.Errorf("StatsOfChangingUseCase - GetAddressWithBiggestChange - getAddressChangeMap: %w", err)
	}
	// Returning result of finding address with biggest changing.
	return uc.getMaxChanging(addresses, currentBlock, int64(countOfLastBlocks)), nil
}

// Get map which store addresses and changes in last countOfLastBlocks blocks.
func (uc *StatsOfChangingUseCase) getAddressChangeMap(
	ctx context.Context,
	currentBlock *big.Int,
	countOfLastBlocks int,
) (map[string]*big.Int, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	var errCh = make(chan error, countOfLastBlocks)

	var mutex sync.Mutex

	addresses := make(map[string]*big.Int, uc.averageAddressCountInBlock*countOfLastBlocks)

	pool := make(chan struct{}, uc.maxGoroutines)

	wg.Add(countOfLastBlocks)

	// Starting from oldest blocks for store earliest blocks
	firstBlock := new(big.Int).Sub(currentBlock, big.NewInt(int64(countOfLastBlocks-1)))

	// Launching goroutines pool
	// where we process addresses with changes to each block and
	// calculate the total value in addresses map.
	for i := 0; i < countOfLastBlocks; i++ {
		pool <- struct{}{}

		go func(offset int) {
			defer wg.Done()
			defer func() { <-pool }()

			blockNumber := new(big.Int).Add(firstBlock, big.NewInt(int64(offset)))

			// Get addresses with changes by blockNumber.
			// If we get error, sending it in errCh and cancelling context.
			chs, err := uc.getAddressWithChanges(ctx, blockNumber)
			if err != nil {
				errCh <- err

				return
			}
			// Adding new change into addresses.
			mutex.Lock()
			for addr, change := range chs {
				if addresses[addr] == nil {
					addresses[addr] = new(big.Int)
				}

				addresses[addr] = new(big.Int).Add(addresses[addr], change)
			}
			mutex.Unlock()
		}(i)

		// If some gouroutine is already ended with error, returning this error
		select {
		case err := <-errCh:
			cancel()
			return nil, err
		default:
		}
	}
	wg.Wait()

	close(pool)
	close(errCh)

	// When all goroutines is done, we checking channel with errors.
	if err := <-errCh; err != nil {
		return nil, err
	}

	return addresses, nil
}

// Getting addresses with changes by number of block.
func (uc *StatsOfChangingUseCase) getAddressWithChanges(
	ctx context.Context,
	blockNumber *big.Int,
) (map[string]*big.Int, error) {
	// Trying to get values from cache
	if cachedResult, ok := uc.cache.Get(blockNumber.String()); ok {
		res, ok := cachedResult.(map[string]*big.Int)
		if ok {
			return res, nil
		}
	}

	chs := make(map[string]*big.Int, uc.averageAddressCountInBlock)

	// Making request to web api
	trs, err := uc.webAPI.GetTransactionsByBlockNumber(ctx, blockNumber)
	if err != nil {
		return nil,
			fmt.Errorf("StatsOfChangingUseCase - getAddressWithChanges - uc.webAPI.GetTransactionsByBlockNumber: %w", err)
	}

	// Calculating amount that the sender spent and receiver got
	for _, t := range trs {
		totalGas := new(big.Int).Mul(t.Gas, t.GasPrice)
		totalFrom := new(big.Int).Add(t.Value, totalGas)

		if chs[t.From] == nil {
			chs[t.From] = new(big.Int)
		}

		if chs[t.To] == nil {
			chs[t.To] = new(big.Int)
		}

		chs[t.From] = new(big.Int).Sub(chs[t.From], totalFrom)
		chs[t.To] = new(big.Int).Add(chs[t.To], t.Value)
	}
	// Adding value in cache
	uc.cache.Add(blockNumber.String(), chs)

	return chs, nil
}

// Getting max change in map with addresses and them changing.
func (uc *StatsOfChangingUseCase) getMaxChanging(
	addresses map[string]*big.Int,
	currentBlock *big.Int,
	countOfLastBlocks int64,
) *entity.BiggestChange {
	maxChange := big.NewInt(0)
	res := &entity.BiggestChange{
		LastBlock:     int2hex(currentBlock),
		CountOfBlocks: countOfLastBlocks,
	}
	// Comparing the current maxChange with current amount
	for addr, amount := range addresses {
		if new(big.Int).Abs(amount).Cmp(new(big.Int).Abs(maxChange)) > 0 {
			res.Address = addr
			maxChange = amount
		}
	}
	// If number is not positive IsRecieved will be false
	if maxChange.Cmp(big.NewInt(0)) > 0 {
		res.IsRecieved = true
	}
	// Amount will be unsigned
	res.Amount = int2hex(new(big.Int).Abs(maxChange))

	return res
}

func int2hex(i *big.Int) string {
	return fmt.Sprintf("%#x", i)
}
