package main

import "context"

type Service struct {
}

type CodeExecInfo struct {
	Lang    string
	Content string
	Args    []string
}

func (s *Service) executeCode(ctx context.Context, info CodeExecInfo) (string, error) {
	return "", nil
}
