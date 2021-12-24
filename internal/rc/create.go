package rc

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// Create creates a container then returns the container id
func (c *Client) Create(ctx context.Context, hostConfig *container.HostConfig) (string, error) {
	ct, err := c.containerPort.ContainerCreate(ctx, c.containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return "", err
	}
	log.Printf("warnings: %v\n", ct.Warnings)

	return ct.ID, c.containerPort.ContainerStart(ctx, ct.ID, types.ContainerStartOptions{})
}
