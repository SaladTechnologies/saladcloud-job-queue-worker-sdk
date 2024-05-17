package main

import (
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/saladtechnologies/saladcloud-job-queue-worker-sdk/internal/loggers"
	"github.com/saladtechnologies/saladcloud-job-queue-worker-sdk/internal/workers"
)

var logger = loggers.Logger

func main() {
	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load()
	}

	var config workers.Config
	if err := env.Parse(&config); err != nil {
		logger.Fatalln(err)
	}

	logger.Infof("starting with the config %+v", config)
	client := workers.New(&config)
	client.Run()
}
