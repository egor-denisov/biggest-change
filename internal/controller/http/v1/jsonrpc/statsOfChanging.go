package jsonrpc

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/egor-denisov/biggest-change/internal/entity"
	"github.com/egor-denisov/biggest-change/internal/usecase"
	sl "github.com/egor-denisov/biggest-change/pkg/logger"
)

var _defaultCountOfBlocks uint = 100

type StatsOfChangingService struct {
	l  *slog.Logger
	sc usecase.StatsOfChanging
}

func NewStatsOfChangingService(l *slog.Logger, sc usecase.StatsOfChanging) *StatsOfChangingService {
	return &StatsOfChangingService{
		l:  l,
		sc: sc,
	}
}

// Getting biggest change for json rpc endpoint.
func (s *StatsOfChangingService) GetBiggestChange(
	r *http.Request,
	args *GetBiggestChangeArgs,
	result *GetBiggestChangeResult,
) error {
	if args.CountOfBlocks == 0 {
		args.CountOfBlocks = _defaultCountOfBlocks
	}

	res, err := s.sc.GetAddressWithBiggestChange(r.Context(), args.CountOfBlocks)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return entity.ErrProcessTimeout
		}

		s.l.Error("jsonrpc - GetBiggestChange", sl.Err(err))

		return entity.ErrInternalServer
	}

	*result = res

	return nil
}

type GetBiggestChangeArgs struct {
	CountOfBlocks uint `json:"countOfBlocks"`
}

type GetBiggestChangeResult *entity.BiggestChange
