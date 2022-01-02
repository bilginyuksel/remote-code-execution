package rc

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
)

func (c *Client) ForceRemove(ctx context.Context, id string) {
	if err := c.containerPort.ContainerStop(ctx, id, nil); err != nil {
		log.Printf("could not stop container, err: %v\n", err)
	}

	if err := c.containerPort.ContainerRemove(ctx, id, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   true,
		Force:         true,
	}); err != nil {
		log.Printf("container remove: %v\n", err)
	}
	log.Println("container removed")
}
