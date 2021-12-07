package main

import (
	"context"
	"log"

	"github.com/docker/docker/client"
	"github.com/yudai/pp"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	manager := NewManager(cli, "ubuntu:latest")
	container, err := manager.Create(ctx)
	if err != nil {
		panic(err)
	}
	log.Println(container)

	type Mock struct {
		Title string `json:"title"`
	}

	res, err := manager.Exec(ctx, container.ID, []string{"echo", `{"title": "MiTitle"}`})
	if err != nil {
		panic(err)
	}
	var mock Mock
	if err := res.Unmarshal(&mock); err != nil {
		panic(err)
	}
	pp.Println(mock)
}
