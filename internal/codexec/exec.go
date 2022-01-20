package codexec

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"
)

type ExecutionRes struct {
	Output        string
	ExecutionTime time.Duration
}

func (c *Codexec) ExecOnce(ctx context.Context, info ExecutionInfo) (*ExecutionRes, error) {
	if !supportedLanguages.IsSupported(info.Lang) {
		return nil, errors.New("language is not supported")
	}

	containerID, err := c.containerClient.Create(ctx, c.hostConfig)
	if err != nil {
		return nil, err
	}
	defer func() { go c.containerClient.ForceRemove(context.Background(), containerID) }()
	return c.exec(ctx, containerID, info)
}

func (c *Codexec) Exec(ctx context.Context, containerID string, info ExecutionInfo) (*ExecutionRes, error) {
	if !supportedLanguages.IsSupported(info.Lang) {
		return nil, errors.New("language is not supported")
	}

	return c.exec(ctx, containerID, info)
}

func (c *Codexec) exec(ctx context.Context, containerID string, info ExecutionInfo) (*ExecutionRes, error) {
	targetFileDir, err := c.findOrWriteContent(supportedLanguages[info.Lang].Filename(), info.Content)
	if err != nil {
		return nil, err
	}
	log.Printf("container: %s, path= %s, cmd= %s\n", containerID, targetFileDir, info.Command())
	res, err := c.containerClient.Exec(ctx, containerID, targetFileDir, info.Command())
	return &ExecutionRes{
		Output:        res.Buffer(),
		ExecutionTime: res.ExecutionTime,
	}, err
}

func (c *Codexec) findOrWriteContent(filename, content string) (string, error) {
	if dir, err := c.cache.Get(content).Result(); err == nil {
		log.Printf("content found in cache, dir: %s\n", dir)
		return dir, nil
	}

	dir, err := c.writeContent(filename, content)
	if err != nil {
		return "", err
	}

	status := c.cache.Set(content, dir, _cacheClearDuration)
	log.Println(status)
	return dir, nil
}

func (c *Codexec) writeContent(filename, content string) (string, error) {
	sourceFilepath, err := c.write(MountSource, filename, content)
	if err != nil {
		return "", err
	}
	targetFilepath := strings.ReplaceAll(sourceFilepath, MountSource, MountTarget)
	targetFileDir := strings.ReplaceAll(targetFilepath, filename, "")
	return targetFileDir, nil
}
