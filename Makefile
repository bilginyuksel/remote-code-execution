mockgen:
	mockgen -destination=mocks/mock_container_port.go -package mocks github.com/codigician/remote-code-execution/internal/rc ContainerPort
	mockgen -destination=mocks/mock_io_readcloser.go -package mocks io ReadCloser

test:
	go test ./... -v

code-coverage:
	go test -coverpkg=./... -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out | grep total | awk '{print $3}'

lint:
	golangci-lint run