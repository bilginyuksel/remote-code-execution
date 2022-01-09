package rc

import (
	"bytes"
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
)

type ExecRes struct {
	Buffer  []byte
	Success bool
}

func NewExecRes(outBuf, errBuf []byte) *ExecRes {
	res := &ExecRes{
		Buffer:  outBuf,
		Success: true,
	}

	if outBuf == nil {
		res.Buffer = errBuf
		res.Success = false
	}

	return res
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

	attachExecRes, err := c.containerPort.ContainerExecAttach(ctx, createExecResponse.ID, types.ExecStartCheck{})
	if err != nil {
		log.Printf("container exec attach: %v\n", err)
		return nil, err
	}

	var outBuffer, errBuffer bytes.Buffer
	_, err = stdcopy.StdCopy(&outBuffer, &errBuffer, attachExecRes.Conn)
	log.Printf("container exec attached, outBuffer: %v, errBuffer: %v\n", outBuffer.String(), errBuffer.String())
	return NewExecRes(outBuffer.Bytes(), errBuffer.Bytes()), err
}
