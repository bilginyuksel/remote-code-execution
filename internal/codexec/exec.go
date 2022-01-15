package codexec

import (
	"context"
	"errors"
	"log"
	"strings"
)

func (c *Codexec) ExecOnce(ctx context.Context, info ExecutionInfo) ([]byte, error) {
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

func (c *Codexec) Exec(ctx context.Context, containerID string, info ExecutionInfo) ([]byte, error) {
	if !supportedLanguages.IsSupported(info.Lang) {
		return nil, errors.New("language is not supported")
	}

	return c.exec(ctx, containerID, info)
}

func (c *Codexec) exec(ctx context.Context, containerID string, info ExecutionInfo) ([]byte, error) {
	sourceFilepath, err := c.write(MountSource, supportedLanguages[info.Lang].Filename(), info.Content)
	if err != nil {
		return nil, err
	}
	targetFilepath := strings.ReplaceAll(sourceFilepath, MountSource, MountTarget)
	targetFileDir := strings.ReplaceAll(targetFilepath, supportedLanguages[info.Lang].Filename(), "")
	log.Printf("container: %s, path= %s, cmd= %s\n", containerID, targetFilepath, info.Command())
	res, err := c.containerClient.Exec(ctx, containerID, targetFileDir, info.Command())
	return res.Buffer, err
}
