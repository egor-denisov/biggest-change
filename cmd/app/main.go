package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/egor-denisov/biggest-change/config"
	app "github.com/egor-denisov/biggest-change/internal/app"
	"github.com/egor-denisov/biggest-change/pkg/logger"
)

func main() {
	// Init configuration
	cfg := config.MustLoad()

	// Init logger
	log := logger.SetupLogger(cfg.Log.Level)

	// Init application
	if cfg.API.Url == "" {
		panic("app cannot be started without url")
	}

	application := app.New(log, cfg)

	// Run server
	go func() {
		application.HTTPServer.MustRun()
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Info("Starting graceful shutdown")

	application.HTTPServer.Stop()

	log.Info("Gracefully stopped")
}
