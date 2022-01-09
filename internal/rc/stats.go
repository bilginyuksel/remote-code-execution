package rc

import (
	"context"

	"github.com/docker/docker/api/types"
)

func (c *Client) Stats(ctx context.Context, id string) (types.ContainerStats, error) {
	return c.containerPort.ContainerStatsOneShot(ctx, id)
}
