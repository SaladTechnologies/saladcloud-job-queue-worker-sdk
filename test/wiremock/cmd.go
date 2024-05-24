package main

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/saladtechnologies/saladcloud-job-queue-worker-sdk/pkg/config"
	"github.com/saladtechnologies/saladcloud-job-queue-worker-sdk/pkg/jobs"
	"github.com/saladtechnologies/saladcloud-job-queue-worker-sdk/pkg/log"
	"github.com/saladtechnologies/saladcloud-job-queue-worker-sdk/pkg/workers"
)

func main() {
	ctx := context.Background()
	logger := slog.Default()
	wiremock, err := StartWiremockContainer(ctx, []WiremockMapping{
		{
			Name:    "status.json",
			Content: "{\"request\":{\"method\":\"GET\",\"urlPath\":\"/v1/status\"},\"response\":{\"status\":200,\"headers\":{\"Content-Type\":\"application/json\"},\"jsonBody\":{\"ready\":true,\"started\":true}}}",
		},
		{
			Name:    "token.json",
			Content: "{\"request\":{\"method\":\"GET\",\"urlPath\":\"/v1/token\"},\"response\":{\"status\":200,\"headers\":{\"Content-Type\":\"application/json\"},\"jsonBody\":{\"jwt\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzYWxhZF9tYWNoaW5lX2lkIjoiMTIzIn0.sDcpi60mMSxwN5ZQcX8oyhgjAjdQBXjubLhOkjNfXKE\"}}}",
		},
	})
	if err != nil {
		logger.Error("failed to start wiremock", "err", err)
		os.Exit(1)
	}
	defer func() {
		err := wiremock.Stop(ctx)
		if err != nil {
			logger.Error("failed to stop wiremock", "err", err)
			os.Exit(1)
		}
	}()

	c, err := config.NewConfigFromEnv()
	if err != nil {
		logger.Error("failed to create config", "err", err)
		os.Exit(1)
	}

	c.MetadataURI = wiremock.BaseURL()
	c.ServiceEndpoint = "localhost:7070"
	c.ServiceUseTLS = false

	w := workers.NewWorker(c, executeJob)
	err = w.Run(log.WithLogger(ctx, logger))
	if err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("failed to run worker", "err", err)
		os.Exit(1)
	}

	logger.Info("successfully ran worker")
	os.Exit(0)
}

func executeJob(ctx context.Context, job jobs.HTTPJob) ([]byte, error) {
	delayDuration := 500*time.Millisecond + time.Duration(rand.Intn(1000))*time.Millisecond
	delay := time.NewTimer(delayDuration)
	select {
	case <-ctx.Done():
		delay.Stop()
		return nil, ctx.Err()
	case <-delay.C:
	}

	output := []byte("{\"status\":\"done\"}")
	return output, nil
}
