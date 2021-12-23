package rc

import (
	"bytes"
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
)

// Exec executes the given command to the container with the given id
func (c *Client) Exec(ctx context.Context, containerID string, cmd []string) ([]byte, error) {
	createExecResponse, err := c.containerPort.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		Cmd:          cmd,
		AttachStderr: true,
		AttachStdout: true,
	})
	if err != nil {
		log.Printf("container exec create: %v, cmd: %v\n", err, cmd)
		return nil, err
	}

	attachExecRes, err := c.containerPort.ContainerExecAttach(ctx, createExecResponse.ID, types.ExecStartCheck{})
	if err != nil {
		log.Printf("container exec attach: %v\n", err)
		return nil, err
	}

	var outBuffer, errBuffer bytes.Buffer
	_, err = stdcopy.StdCopy(&outBuffer, &errBuffer, attachExecRes.Conn)

	return outBuffer.Bytes(), err
}
