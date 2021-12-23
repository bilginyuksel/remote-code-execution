package rc

import (
	"context"

	"github.com/docker/docker/api/types"
)

// List the containers
func (c *Client) List(ctx context.Context) ([]string, error) {
	containers, err := c.containerPort.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	var containerIDs []string
	for idx := range containers {
		containerIDs = append(containerIDs, containers[idx].ID)
	}

	return containerIDs, nil
}
