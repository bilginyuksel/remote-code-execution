mockgen:
	mockgen -destination=internal/mocks/mock_container_port.go -package mocks github.com/codigician/remote-code-execution/internal/rc ContainerPort
	mockgen -destination=internal/mocks/mock_io_readcloser.go -package mocks io ReadCloser
	mockgen -destination=internal/mocks/mock_conn.go -package mocks net Conn
	mockgen -destination=internal/mocks/mock_container_client.go -package mocks github.com/codigician/remote-code-execution/internal/codexec ContainerClient
	mockgen -destination=internal/mocks/mock_codexecutor.go -package mocks github.com/codigician/remote-code-execution/internal/codexec Codexecutor 
	mockgen -destination=internal/mocks/mock_remote_codexecutor_service.go -package mocks github.com/codigician/remote-code-execution/internal/handler RemoteCodeExecutorService

unit-test:
	go test ./... -v -short

test:
	go test ./... -v

code-coverage:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out | grep total | awk '{print $3}'

lint:
	golangci-lint run

up:
	chmod +x ./scripts/run_docker.sh
	./scripts/run_docker.sh