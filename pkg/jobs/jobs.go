package jobs

import "context"

// Represents an HTTP request job from a SaladCloud job queue.
type HTTPJob struct {
	// The identifier of the job.
	JobId string
	// The HTTP request method.
	RequestMethod string
	// The HTTP request URL.
	RequestURL string
	// The HTTP request body.
	RequestBody []byte
}

// Represents a function that executes an HTTP request job from a SaladCloud job
// queue. The returned value is the HTTP response body.
type HTTPJobExecutor = func(ctx context.Context, job HTTPJob) ([]byte, error)
