mockgen:
	mockgen -destination=internal/mocks/mock_container_port.go -package mocks github.com/codigician/remote-code-execution/internal/rc ContainerPort
	mockgen -destination=internal/mocks/mock_io_readcloser.go -package mocks io ReadCloser
	mockgen -destination=internal/mocks/mock_conn.go -package mocks net Conn
	mockgen -destination=internal/mocks/mock_container_client.go -package mocks 

test:
	go test ./... -v

code-coverage:
	go test -coverpkg=./... -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out | grep total | awk '{print $3}'

lint:
	golangci-lint run