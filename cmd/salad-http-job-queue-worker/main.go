package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/saladtechnologies/saladcloud-job-queue-worker-sdk/pkg/config"
	"github.com/saladtechnologies/saladcloud-job-queue-worker-sdk/pkg/jobs"
	"github.com/saladtechnologies/saladcloud-job-queue-worker-sdk/pkg/log"
	"github.com/saladtechnologies/saladcloud-job-queue-worker-sdk/pkg/workers"
)

func main() {
	defaultLogger := slog.Default()
	c, err := config.NewConfigFromEnv()
	if err != nil {
		defaultLogger.Error("failed to load environment variables", "err", err)
		os.Exit(1)
	}

	var leveler slog.Leveler
	switch c.LogLevel {
	case "debug":
		leveler = slog.LevelDebug
	case "info":
		leveler = slog.LevelInfo
	case "warn":
		leveler = slog.LevelWarn
	case "error":
		leveler = slog.LevelError
	default:
		defaultLogger.Error("invalid SALAD_LOG_LEVEL environment variable value", "level", c.LogLevel)
		os.Exit(1)
	}

	logger := slog.New(log.NewLeveledHandler(defaultLogger.Handler(), leveler))
	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		logger.Info("stopping SaladCloud Job Queue worker...")
		cancel()
	}()

	logger.Info("starting SaladCloud Job Queue worker...")
	w := workers.NewWorker(c, executeJob)
	err = w.Run(log.WithLogger(ctx, logger))
	if err != nil && !errors.Is(err, context.Canceled) {
		defaultLogger.Error("failed to run SaladCloud Job Queue worker", "err", err)
		os.Exit(1)
	}

	logger.Info("stopped SaladCloud Job Queue worker")
	os.Exit(0)
}

func executeJob(ctx context.Context, job jobs.HTTPJob) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, job.RequestMethod, job.RequestURL, bytes.NewReader(job.RequestBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	output, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return output, nil
}
