package main

import (
	"context"
	"fmt"
	"log"
)

type Service struct {
	manager *ContainerManager
}

type CodeExecInfo struct {
	Lang    string
	Content string
	Args    []string
}

func (s *Service) executeCode(ctx context.Context, info CodeExecInfo) (string, error) {
	containerID, err := s.manager.Create(ctx)
	if err != nil {
		return "", err
	}
	// containerID := "f2e813f33164"

	if _, err := s.manager.Exec(ctx, containerID, []string{"touch", "main.go"}); err != nil {
		log.Println("main.go file could not created")
		return "", err
	}

	if _, err := s.manager.Exec(ctx, containerID, []string{"echo", fmt.Sprintf("'%s'", info.Content), ">", "main.go"}); err != nil {
		log.Println("could not echo the content to main.go")
		return "", err
	}

	baseRunner := []string{"/usr/local/go/bin/go", "run", "app/main.go"}
	baseRunner = append(baseRunner, info.Args...)
	response, err := s.manager.Exec(ctx, containerID, baseRunner)
	if err != nil {
		log.Println("could not run main.go")
		return "", err
	}
	return string(response), err
}
