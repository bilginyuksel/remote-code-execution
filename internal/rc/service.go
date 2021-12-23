package rc

import (
	"context"
	"fmt"
	"log"
)

type Service struct {
	client *Client
}

type CodeExecInfo struct {
	Lang    string
	Content string
	Args    []string
}

func (s *Service) executeCode(ctx context.Context, info CodeExecInfo) (string, error) {
	containerID, err := s.client.Create(ctx, nil)
	if err != nil {
		return "", err
	}
	// containerID := "f2e813f33164"

	if _, err := s.client.Exec(ctx, containerID, []string{"touch", "main.go"}); err != nil {
		log.Println("main.go file could not created")
		return "", err
	}

	if _, err := s.client.Exec(ctx, containerID, []string{"echo", fmt.Sprintf("'%s'", info.Content), ">", "main.go"}); err != nil {
		log.Println("could not echo the content to main.go")
		return "", err
	}

	baseRunner := []string{"/usr/local/go/bin/go", "run", "app/main.go"}
	baseRunner = append(baseRunner, info.Args...)
	response, err := s.client.Exec(ctx, containerID, baseRunner)
	if err != nil {
		log.Println("could not run main.go")
		return "", err
	}
	return string(response), err
}
