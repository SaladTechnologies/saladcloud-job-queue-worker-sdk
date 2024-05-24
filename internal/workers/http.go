package workers

import (
	"context"
	"log/slog"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	metadata_api "github.com/saladtechnologies/salad-cloud-job-queue-worker/internal/apis/metadata"
	"github.com/saladtechnologies/salad-cloud-job-queue-worker/pkg/log"
)

// Creates a SaladCloud Instance Metadata Service (IMDS) client.
func newMetadataClient(ctx context.Context, addr string) (*metadata_api.ClientWithResponses, error) {
	rc := retryablehttp.NewClient()
	rc.Logger = &leveledLogger{logger: log.FromContext(ctx)}
	rc.RetryMax = 5
	rc.RetryWaitMin = 1 * time.Second
	rc.RetryWaitMax = 2 * time.Minute

	c, err := metadata_api.NewClientWithResponses(addr, metadata_api.WithHTTPClient(rc.StandardClient()))
	if err != nil {
		return nil, err
	}

	return c, nil
}

type leveledLogger struct {
	logger *slog.Logger
}

func (l *leveledLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Debug(msg, keysAndValues...)
}

func (l *leveledLogger) Error(msg string, keysAndValues ...interface{}) {
	l.logger.Error(msg, keysAndValues...)
}

func (l *leveledLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Info(msg, keysAndValues...)
}

func (l *leveledLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warn(msg, keysAndValues...)
}
