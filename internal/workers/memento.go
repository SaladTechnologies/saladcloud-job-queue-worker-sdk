package workers

import "github.com/saladtechnologies/saladcloud-job-queue-worker-sdk/pkg/gen"

type memento struct {
	job        *gen.Job
	output     []byte
	completion bool
}

func (m *memento) clear() {
	m.job = nil
	m.completion = false
}

func (m *memento) rememberRejection(job *gen.Job) {
	m.clear()
	m.job = job
}

func (m *memento) rememberCompletion(job *gen.Job, output []byte) {
	m.clear()
	m.job = job
	m.output = output
	m.completion = true
}
