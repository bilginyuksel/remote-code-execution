package rc

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// List the running containers which are created by the image
func (c *Client) List(ctx context.Context) ([]string, error) {
	containers, err := c.containerPort.ContainerList(ctx,
		types.ContainerListOptions{
			Filters: filters.NewArgs(
				filters.KeyValuePair{
					Key:   "status",
					Value: "running",
				},
				filters.KeyValuePair{
					Key:   "ancestor",
					Value: c.containerConfig.Image,
				}),
		})

	var containerIDs []string
	for idx := range containers {
		containerIDs = append(containerIDs, containers[idx].ID)
	}

	return containerIDs, err
}
