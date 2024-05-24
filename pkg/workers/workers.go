package workers

import (
	"context"

	"github.com/saladtechnologies/salad-cloud-job-queue-worker/internal/workers"
	"github.com/saladtechnologies/salad-cloud-job-queue-worker/pkg/config"
	"github.com/saladtechnologies/salad-cloud-job-queue-worker/pkg/jobs"
)

// Represents a worker bound to a SaladCloud job queue.
type Worker struct {
	config   config.Config
	executor jobs.HTTPJobExecutor
}

// Creates a new worker bound to a SaladCloud job queue. The worker will use the
// given function to execute jobs as HTTP requests.
func NewWorker(config config.Config, executor jobs.HTTPJobExecutor) *Worker {
	return &Worker{
		config:   config,
		executor: executor,
	}
}

// Runs the worker until the given context is cancelled.
func (w *Worker) Run(ctx context.Context) error {
	return workers.Run(ctx, w.config, w.executor)
}
