package main

import (
	"context"
	"log"

	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// execResponse, err := cli.ContainerExecCreate(context.Background(), res.ID, types.ExecConfig{
	// 	Cmd: []string{"echo", "'hello world'"},
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// execAttach, err := cli.ContainerExecAttach(ctx, execResponse.ID, types.ExecStartCheck{
	// 	Detach: true,
	// 	Tty:    true,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// stdcopy.StdCopy(os.Stdout, os.Stderr, execAttach.Reader)

	// stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	ctx := context.Background()
	manager := NewManager(cli, "ubuntu:latest")
	container, err := manager.Create(ctx)
	if err != nil {
		panic(err)
	}
	log.Println(container)

	res, err := manager.Exec(ctx, container.ID, []string{"echo", "hello world"})
	if err != nil {
		panic(err)
	}
	log.Println(res)
}
