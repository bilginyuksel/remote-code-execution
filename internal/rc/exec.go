package rc

import (
	"bytes"
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
)

// Exec executes the given command to the container with the given id
func (c *Client) Exec(ctx context.Context, containerID, workingDir string, cmd []string) ([]byte, error) {
	createExecResponse, err := c.containerPort.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		Cmd:          cmd,
		AttachStderr: true,
		AttachStdout: true,
		WorkingDir:   workingDir,
	})
	if err != nil {
		log.Printf("container exec create: %v, cmd: %v\n", err, cmd)
		return nil, err
	}
	log.Printf("container exec created, id: %s\n", createExecResponse.ID)

	attachExecRes, err := c.containerPort.ContainerExecAttach(ctx, createExecResponse.ID, types.ExecStartCheck{})
	if err != nil {
		log.Printf("container exec attach: %v\n", err)
		return nil, err
	}
	log.Println("container exec attached")

	var outBuffer, errBuffer bytes.Buffer
	_, err = stdcopy.StdCopy(&outBuffer, &errBuffer, attachExecRes.Conn)

	return outBuffer.Bytes(), err
}
