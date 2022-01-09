package codexec

import (
	"context"
	"fmt"
	"os"

	"github.com/codigician/remote-code-execution/internal/rc"
	"github.com/docker/docker/api/types/container"
	"github.com/google/uuid"
)

const (
	MountSource = "target"
	MountTarget = "/app"
)

type (
	ContainerClient interface {
		Create(ctx context.Context, hostConfig *container.HostConfig) (string, error)
		Exec(ctx context.Context, id, workingDir string, cmd []string) (*rc.ExecRes, error)
		ForceRemove(ctx context.Context, id string)
	}

	Write func(baseDir, filename, content string) (string, error)

	Codexec struct {
		containerClient ContainerClient
		hostConfig      *container.HostConfig
		write           Write
	}
)

func New(containerClient ContainerClient, hostConfig *container.HostConfig, write Write) *Codexec {
	return &Codexec{
		containerClient: containerClient,
		hostConfig:      hostConfig,
		write:           write,
	}
}

func WriteFile(baseDir, filename, content string) (string, error) {
	const permission = 0777

	directory := fmt.Sprintf("%s/%s", baseDir, uuid.NewString())
	filepath := fmt.Sprintf("%s/%s", directory, filename)
	// if mkdir returns error also WriteFile will throw an error so no need to handle mkdir error
	_ = os.Mkdir(directory, permission)
	return filepath, os.WriteFile(filepath, []byte(content), permission)
}
