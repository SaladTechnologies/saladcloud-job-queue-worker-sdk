
build: clean protos
	GOOS=linux go build -o salad-queue-client cmd/main.go

protos:
	mkdir -p pkg/gen
	protoc ../../../proto/packages/saladcloud-job-queue-worker/salad/grpc/saladcloud_job_queue_worker/v1alpha/*.proto \
		--go_out=pkg/gen \
		--go_opt=paths=source_relative \
		--go_opt=Mservice_v1alpha.proto=salad.com/qworker/gen \
		--go-grpc_out=pkg/gen \
		--go-grpc_opt=Mservice_v1alpha.proto=salad.com/qworker/gen \
		--go-grpc_opt=paths=source_relative \
		--proto_path=../../../proto/packages/saladcloud-job-queue-worker/salad/grpc/saladcloud_job_queue_worker/v1alpha/

clean:
	rm -rf pkg/gen salad-queue-client

lint:
	golangci-lint run ./...

