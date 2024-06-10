package workers

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"

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
		Version:  VersionStr(),
	}
}

// Runs the worker until the given context is cancelled.
func (w *Worker) Run(ctx context.Context) error {
	return workers.Run(ctx, w.config, w.executor, w.Version)
}

// Formats a complete version string
func VersionStr() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Sprintf("%s", strings.TrimPrefix(Version, "v"))
	}
	var commit string
	var modified string
	for _, kv := range bi.Settings {
		if kv.Key == "vcs.revision" {
			commit = kv.Value
			if len(commit) > 7 {
				commit = commit[:7]
			}
		}
		if kv.Key == "vcs.modified" {
			modified = "-dirty"
		}
	}
	return fmt.Sprintf("%s+%s%s", strings.TrimPrefix(Version, "v"), commit, modified)
}
