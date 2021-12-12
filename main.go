package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	manager := NewContainerManager(cli, &container.Config{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"bash"},
		Image:        "ubuntu:latest",
	})
	containerID, err := manager.Create(ctx)
	if err != nil {
		panic(err)
	}
	log.Println(containerID)

	startTime := time.Now()
	group := sync.WaitGroup{}
	group.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			res, err := manager.Exec(ctx, containerID, []string{"echo", `{"title": "MiTitle"}`})
			if err != nil {
				panic(err)
			}
			log.Println(string(res))
			group.Done()
		}()
	}

	group.Wait()
	log.Printf("time elapsed: %v\n", time.Since(startTime))
}
