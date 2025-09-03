package main

import (
	"log"

	"github.com/Prabhat7saini/Basic-Setup/cmd/app"
	"github.com/Prabhat7saini/Basic-Setup/config"
	"github.com/Prabhat7saini/Basic-Setup/shared/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := logger.Init(cfg.Log)

	app := app.NewApp(cfg, logger)
	app.Run()
}
