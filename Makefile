
build: clean
	GOOS=linux go build -o salad-job-queue-worker cmd/main.go

clean:
	rm -f salad-job-queue-worker

lint:
	golangci-lint run ./...

