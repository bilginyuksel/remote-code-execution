package main

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/yudai/pp"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(ctx, "docker.io/library/ubuntu", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	res, err := cli.ContainerCreate(ctx, &container.Config{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"bash"},
		Image:        "ubuntu",
		Volumes:      map[string]struct{}{},
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, res.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	// statusChannel, errChannel := cli.ContainerWait(ctx, res.ID, container.WaitConditionNotRunning)
	// select {
	// case err := <-errChannel:
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// case <-statusChannel:
	// 	fmt.Println("container status, ", statusChannel)
	// }

	out, err := cli.ContainerLogs(ctx, res.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	execResponse, err := cli.ContainerExecCreate(ctx, res.ID, types.ExecConfig{
		Cmd: []string{"echo", "'hello world'"},
	})
	if err != nil {
		panic(err)
	}
	execAttach, err := cli.ContainerExecAttach(ctx, execResponse.ID, types.ExecStartCheck{
		Detach: true,
		Tty:    true,
	})
	if err != nil {
		panic(err)
	}
	stdcopy.StdCopy(os.Stdout, os.Stderr, execAttach.Reader)

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		pp.Println(container)
	}

	// manager := NewManager(cli, "docker.io/library/ubuntu")
	// container, err := manager.Create(context.Background())
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(container)
}
