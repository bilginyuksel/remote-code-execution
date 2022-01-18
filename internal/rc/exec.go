package rc

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
)

type ExecRes struct {
	outBuffer     string
	errBuffer     string
	ExecutionTime time.Duration
}

func NewExecRes(outBuffer, errBuffer string, executionTime time.Duration) *ExecRes {
	return &ExecRes{
		outBuffer:     outBuffer,
		errBuffer:     errBuffer,
		ExecutionTime: executionTime,
	}
}

func (e *ExecRes) Buffer() string {
	if e.outBuffer != "" {
		return e.outBuffer
	}
	return e.errBuffer
}

// Exec executes the given command to the container with the given id
// returns outBuffer, errBuffer, error
func (c *Client) Exec(ctx context.Context, containerID, workingDir string, cmd []string) (*ExecRes, error) {
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

	executionStartTime := time.Now()
	attachExecRes, err := c.containerPort.ContainerExecAttach(ctx, createExecResponse.ID, types.ExecStartCheck{})
	if err != nil {
		log.Printf("container exec attach: %v\n", err)
		return nil, err
	}

	var outBuffer, errBuffer bytes.Buffer
	_, err = stdcopy.StdCopy(&outBuffer, &errBuffer, attachExecRes.Conn)
	log.Printf("container exec attached, outBuffer: %v, errBuffer: %v\n", outBuffer.String(), errBuffer.String())
	return &ExecRes{
		outBuffer:     outBuffer.String(),
		errBuffer:     errBuffer.String(),
		ExecutionTime: time.Since(executionStartTime),
	}, err
}
