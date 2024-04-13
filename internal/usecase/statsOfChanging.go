package usecase

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/egor-denisov/biggest-change/internal/entity"
	lru "github.com/hashicorp/golang-lru"
)

var _maxGoroutines = 50
var _averageAddressCountInBlock = 200
var _cacheSize = 100 // Count of blocks for which the transaction value will be cached

type StatsOfChangingUseCase struct {
	webAPI StatsOfChangingWebAPI
	cache  *lru.Cache
}

func New(w StatsOfChangingWebAPI) *StatsOfChangingUseCase {
	cache, _ := lru.New(_cacheSize)

	return &StatsOfChangingUseCase{
		webAPI: w,
		cache:  cache,
	}
}

// Get address with biggest change in last countOfLastBlocks blocks.
func (w *StatsOfChangingUseCase) GetAddressWithBiggestChange(
	ctx context.Context,
	countOfLastBlocks uint,
) (*entity.BiggestChange, error) {
	// Getting current number of block.
	currentBlock, err := w.webAPI.GetCurrentBlockNumber(ctx)
	if err != nil {
		return nil,
			fmt.Errorf("StatsOfChangingUseCase - GetAddressWithBiggestChange - w.webAPI.GetCurrentBlockNumber: %w", err)
	}
	// Getting map which store addresses and changes in last countOfLastBlocks blocks.
	addresses, err := w.getAddressChangeMap(ctx, currentBlock, int(countOfLastBlocks))
	if err != nil {
		return nil,
			fmt.Errorf("StatsOfChangingUseCase - GetAddressWithBiggestChange - getAddressChangeMap: %w", err)
	}
	// Returning result of finding address with biggest changing.
	return w.getMaxChanging(addresses, currentBlock, int64(countOfLastBlocks)), nil
}

// Get map which store addresses and changes in last countOfLastBlocks blocks.
func (w *StatsOfChangingUseCase) getAddressChangeMap(
	ctx context.Context,
	currentBlock *big.Int,
	countOfLastBlocks int,
) (map[string]*big.Int, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	var errCh = make(chan error, countOfLastBlocks)

	var mutex sync.Mutex

	addresses := make(map[string]*big.Int, _averageAddressCountInBlock*countOfLastBlocks)

	pool := make(chan struct{}, _maxGoroutines)

	// ticker := time.NewTicker(20000 * time.Microsecond)
	// defer ticker.Stop()

	wg.Add(countOfLastBlocks)

	// Starting from oldest blocks for store earliest blocks
	firstBlock := new(big.Int).Sub(currentBlock, big.NewInt(int64(countOfLastBlocks-1)))
	// Launching goroutines pool
	// where we process addresses with changes to each block and
	// calculate the total value in addresses map.
	for i := 0; i < countOfLastBlocks; i++ {
		// <-ticker.C
		pool <- struct{}{}

		go func(offset int) {
			defer wg.Done()
			defer func() { <-pool }()

			blockNumber := new(big.Int).Add(firstBlock, big.NewInt(int64(offset)))

			// Get addresses with changes by blockNumber.
			// If we get error, sending it in errCh and cancelling context.
			chs, err := w.getAddressWithChanges(ctx, blockNumber)
			if err != nil {
				errCh <- err

				cancel()

				return
			}
			// Adding new change into addresses.
			for addr, change := range chs {
				mutex.Lock()
				if addresses[addr] == nil {
					addresses[addr] = new(big.Int)
				}

				addresses[addr] = new(big.Int).Add(addresses[addr], change)
				mutex.Unlock()
			}
		}(i)
	}
	wg.Wait()

	// close(pool)
	close(errCh)

	// When all goroutines is done, we checking channel with errors.
	if err := <-errCh; err != nil {
		return nil, err
	}

	return addresses, nil
}

// Getting addresses with changes by number of block.
func (w *StatsOfChangingUseCase) getAddressWithChanges(
	ctx context.Context,
	blockNumber *big.Int,
) (map[string]*big.Int, error) {
	// Trying to get values from cache
	if cachedResult, ok := w.cache.Get(blockNumber.String()); ok {
		res, ok := cachedResult.(map[string]*big.Int)
		if ok {
			return res, nil
		}
	}

	chs := make(map[string]*big.Int, _averageAddressCountInBlock)

	// Making request to web api
	trs, err := w.webAPI.GetTransactionsByBlockNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
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
	w.cache.Add(blockNumber.String(), chs)

	return chs, nil
}

// Getting max change in map with addresses and them changing.
func (w *StatsOfChangingUseCase) getMaxChanging(
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
