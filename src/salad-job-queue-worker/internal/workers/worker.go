package workers

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"salad.com/qworker/internal/dtos"
	"salad.com/qworker/internal/loggers"
	"salad.com/qworker/pkg/gen"
)

// TODO, clean this up and simplify

var (
	logger = loggers.Logger
)

type Config struct {
	LogLevel        string `env:"LOG_LEVEL" envDefault:"info"`
	ServiceEndpoint string `env:"SERVICE_ENDPOINT" envDefault:"job-queue-worker-api.salad.com:443"`
	MetadataURI     string `env:"SALAD_METADATA_URI" envDefault:"http://169.254.169.254:80"`
	UseTLS          string `env:"USE_TLS" envDefault:"1"`
}

type Worker struct {
	cfg     *Config
	memento memento
	conn    *grpc.ClientConn
	client  gen.JobQueueWorkerServiceClient
	stream  gen.JobQueueWorkerService_AcceptJobsClient
	token   string // TODO associate with the job
}

func New(cfg *Config) *Worker {
	cfg.setLogging()
	return &Worker{
		cfg: cfg,
	}
}

func (w *Worker) Run() {
	for {
		err := w.connectWithBackoff()
		if err != nil {
			logger.Fatalln(err)
		}
		err = w.processMemento()
		if err != nil {
			logger.Errorln(err)
			w.conn.Close()
			continue
		}

		w.handleStream() // exits only if the stream is broken, gracefully or abruptly

		w.conn.Close()
	}
}

func (w *Worker) executeJob(job *gen.Job) error {
	logger.Infof("Executing job %s", job.JobId)
	url := fmt.Sprintf("http://localhost:%d/%s", job.Port, job.Path)
	payload := bytes.NewReader(job.Input)
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	defer resp.Body.Close()
	output, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorln(err)
		if err := w.rejectJob(job); err != nil {
			return err
		}
	}
	return w.completeJob(job, output)
}

func (w *Worker) completeJob(job *gen.Job, output []byte) error {
	logger.Infof("completing job %s", job.JobId)
	w.memento.rememberCompletion(job, output)
	var req gen.CompleteJobRequest
	req.JobId = job.JobId
	req.Output = output
	ctx, err := w.getAuthorizedContext()
	if err != nil {
		logger.Fatalln(err)
	}
	_, err = w.client.CompleteJob(ctx, &req)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	w.memento.clear()
	return nil
}

func (w *Worker) rejectJob(job *gen.Job) error {
	logger.Warnf("rejecting job %s", job.JobId)
	w.memento.rememberRejection(job)
	var req gen.RejectJobRequest
	req.JobId = job.JobId
	ctx, err := w.getAuthorizedContext()
	if err != nil {
		logger.Fatalln(err)
	}
	_, err = w.client.RejectJob(ctx, &req)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	w.memento.clear()
	return nil
}

func (w *Worker) handleStream() {
	for {
		resp, err := w.stream.Recv()
		if err == io.EOF {
			logger.Warningln("end of server stream")
			break
		}
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				logger.Errorf("non-gRPC error: %s", err)
			}
			if st.Code() == codes.Unavailable {
				logger.Warnln("server connection broken")
				break
			}
		}
		if resp == nil {
			logger.Warningln("nil response")
			continue
		}
		if resp.Message == nil {
			logger.Warningln("nil message")
			continue
		}
		switch msg := resp.Message.(type) {
		case *gen.AcceptJobsResponse_Heartbeat:
			logger.Infoln("heartbeat")
			continue
		case *gen.AcceptJobsResponse_Job:
			err := w.executeJob(msg.Job)
			if err != nil {
				break
			}
		}
	}
}

func (cfg *Config) setLogging() {
	for _, level := range logrus.AllLevels {
		if cfg.LogLevel == level.String() {
			logrus.SetLevel(level)
			break
		}
	}
}

func (w *Worker) connectWithBackoff() error {
	tlsConfig := tls.Config{
		InsecureSkipVerify: false,
	}
	var creds credentials.TransportCredentials
	if w.cfg.UseTLS == "" {
		creds = insecure.NewCredentials()
	} else {
		creds = credentials.NewTLS(&tlsConfig)
	}
	var req gen.AcceptJobsRequest

	// keep connecting ...
	sleepMultiplier := 1
	for {
		conn, err := grpc.Dial(w.cfg.ServiceEndpoint, grpc.WithTransportCredentials(creds))
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				return err
			}
			if st.Code() == codes.Unavailable {
				logger.Infof("can't connect, retrying in %d second(s)", sleepMultiplier)
				time.Sleep(time.Second * time.Duration(sleepMultiplier))
				sleepMultiplier *= 2
				continue
			}
		}
		w.conn = conn
		w.client = gen.NewJobQueueWorkerServiceClient(conn)

		ctx, err := w.getAuthorizedContext()
		if err != nil {
			logger.Fatalln(err)
		}

		stream, err := w.client.AcceptJobs(ctx, &req)
		if err != nil {
			logger.Warnln(err)
			st, ok := status.FromError(err)
			if !ok {
				return err
			}
			if st.Code() == codes.Unavailable {
				logger.Warnf("can't connect, retrying in %d second(s)", sleepMultiplier)
				w.conn.Close()
				w.conn = nil
				w.client = nil
				time.Sleep(time.Second * time.Duration(sleepMultiplier))
				sleepMultiplier *= 2
				continue
			}
		}
		w.stream = stream
		return nil
	}
}

func (w *Worker) processMemento() error {
	if w.memento.job == nil {
		return nil
	}
	if w.memento.completion {
		return w.completeJob(w.memento.job, w.memento.output)
	}
	return w.rejectJob(w.memento.job)
}

func (w *Worker) fetchToken() (string, error) {
	if w.token != "" {
		return w.token, nil
	}
	url := fmt.Sprintf("%s/v1/token", w.cfg.MetadataURI)
	logger.Debugf("Fetching token from %s", url)
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected staus code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var token dtos.TokenResponse
	if err = json.Unmarshal(payload, &token); err != nil {
		return "", err
	}
	logger.Tracef("token: %s", token.JWT)
	w.token = token.JWT
	return w.token, nil
}

func (w *Worker) getAuthorizedContext() (context.Context, error) {
	ctx := context.Background()
	token, err := w.fetchToken()
	if err != nil {
		return ctx, err
	}

	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", token))
	return metadata.NewOutgoingContext(ctx, md), nil
}
