// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	"math/big"

	"github.com/egor-denisov/biggest-change/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type (
	StatsOfChanging interface {
		GetAddressWithBiggestChange(ctx context.Context, countOfLastBlocks uint) (*entity.BiggestChange, error)
	}

	StatsOfChangingWebAPI interface {
		GetTransactionsByBlockNumber(ctx context.Context, blockNumber *big.Int) ([]*entity.Transaction, error)
		GetCurrentBlockNumber(ctx context.Context) (*big.Int, error)
	}
)
