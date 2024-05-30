package workers

import (
	"context"
	"fmt"
	"strings"

	"github.com/carlmjohnson/versioninfo"
	"github.com/saladtechnologies/salad-cloud-job-queue-worker/internal/workers"
	"github.com/saladtechnologies/salad-cloud-job-queue-worker/pkg/config"
	"github.com/saladtechnologies/salad-cloud-job-queue-worker/pkg/jobs"
)

var Version = "v0.0.0"

// Represents a worker bound to a SaladCloud job queue.
type Worker struct {
	config   config.Config
	executor jobs.HTTPJobExecutor
	Version  string
}

// Creates a new worker bound to a SaladCloud job queue. The worker will use the
// given function to execute jobs as HTTP requests.
func NewWorker(config config.Config, executor jobs.HTTPJobExecutor) *Worker {
	return &Worker{
		config:   config,
		executor: executor,
		Version:  fmt.Sprintf("%s+%s", strings.TrimPrefix(Version, "v"), versioninfo.Short()),
	}
}

// Runs the worker until the given context is cancelled.
func (w *Worker) Run(ctx context.Context) error {
	return workers.Run(ctx, w.config, w.executor)
}
