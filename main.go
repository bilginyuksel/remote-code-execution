package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/yudai/pp"
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
	container, err := manager.Create(ctx)
	if err != nil {
		panic(err)
	}
	log.Println(container)

	type Mock struct {
		Title string `json:"title"`
	}

	startTime := time.Now()
	group := sync.WaitGroup{}
	group.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			res, err := manager.Exec(ctx, container.ID, []string{"echo", `{"title": "MiTitle"}`})
			if err != nil {
				panic(err)
			}
			var mock Mock
			if err := res.Unmarshal(&mock); err != nil {
				panic(err)
			}
			pp.Println(mock)
			group.Done()
		}()
	}

	group.Wait()
	log.Printf("time elapsed: %v\n", time.Since(startTime))
}
