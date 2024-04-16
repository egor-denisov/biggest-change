package app

import (
	"log/slog"

	"github.com/egor-denisov/biggest-change/config"
	v1 "github.com/egor-denisov/biggest-change/internal/controller/http/v1"
	"github.com/egor-denisov/biggest-change/internal/usecase"
	webapi "github.com/egor-denisov/biggest-change/internal/webapi/getblock"
	"github.com/egor-denisov/biggest-change/pkg/httpserver"

	"github.com/gin-gonic/gin"
)

type App struct {
	HTTPServer *httpserver.Server
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	// Web api
	api := webapi.New(
		cfg.API.URL,
		webapi.RequestCountRPS(cfg.API.Rps),
		webapi.TimeWindowRPS(cfg.API.TimeWindowRPS),
		webapi.Timeout(cfg.API.Timeout),
		webapi.MaxRetries(cfg.API.MaxRetries),
		webapi.TimeBetweenRetries(cfg.API.TimeBetweenRetries),
	)

	// Use case
	statsOfChangingUseCase := usecase.New(
		api,
		usecase.CacheSize(cfg.App.CacheSize),
		usecase.MaxGoroutines(cfg.App.MaxGoroutines),
		usecase.AverageAddressCountInBlock(cfg.App.AverageAddressesInBlock),
		usecase.CountOfBlocks(cfg.App.CountOfBlocks),
	)

	// Init http server
	handler := gin.New()
	v1.NewRouter(handler, log, statsOfChangingUseCase)
	httpServer := httpserver.New(log, handler, httpserver.Port(cfg.HTTP.Port), httpserver.WriteTimeout(cfg.HTTP.Timeout))

	return &App{
		HTTPServer: httpServer,
	}
}
