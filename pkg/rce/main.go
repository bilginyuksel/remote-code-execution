package rce

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	//	tar, err := archive.TarWithOptions("custom-ubuntu/", &archive.TarOptions{})
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	res, err := cli.ImageBuild(ctx, tar, types.ImageBuildOptions{
	//		Dockerfile: "./custom-ubuntu-dockerfile",
	//		Tags:       []string{"custom-ubuntu"},
	//	})
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Println(res)

	manager := NewContainerManager(cli, &container.Config{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"bin/sh"},
		Image:        "custom-ubuntu:latest",
	})

	// //	containerID, err := manager.Create(ctx)
	// //	if err != nil {
	// //		panic(err)
	// //	}
	// //	log.Println(containerID)
	// //
	// //	startTime := time.Now()
	// //	group := sync.WaitGroup{}
	// //	group.Add(100)
	// //	for i := 0; i < 100; i++ {
	// //		go func() {
	// //			res, err := manager.Exec(ctx, containerID, []string{"echo", `{"title": "MiTitle"}`})
	// //			if err != nil {
	// //				panic(err)
	// //			}
	// //			log.Println(string(res))
	// //			group.Done()
	// //		}()
	// //	}
	// //
	// //	group.Wait()
	// //	log.Printf("time elapsed: %v\n", time.Since(startTime))
	service := &Service{manager}
	info := CodeExecInfo{
		Lang: "Golang",
		Content: `package main
		import "fmt"

		func main() {
			fmt.Println("hello world")
		}`,
		Args: []string{"--debug", "life"},
	}
	log.Println(info.Content)

	res, err := service.executeCode(ctx, info)
	if err != nil {
		panic(err)
	}

	log.Println(res)
}
