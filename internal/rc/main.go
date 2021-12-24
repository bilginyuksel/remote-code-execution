package rc

// func main() {
// 	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 	if err != nil {
// 		panic(err)
// 	}

// 	ctx := context.Background()

// 	//	tar, err := archive.TarWithOptions("custom-ubuntu/", &archive.TarOptions{})
// 	//	if err != nil {
// 	//		log.Fatal(err)
// 	//	}
// 	//
// 	//	res, err := cli.ImageBuild(ctx, tar, types.ImageBuildOptions{
// 	//		Dockerfile: "./custom-ubuntu-dockerfile",
// 	//		Tags:       []string{"custom-ubuntu"},
// 	//	})
// 	//	if err != nil {
// 	//		log.Fatal(err)
// 	//	}
// 	//	log.Println(res)

// 	manager := NewClient(cli, &container.Config{
// 		AttachStdin:  true,
// 		AttachStdout: true,
// 		AttachStderr: true,
// 		Tty:          true,
// 		Cmd:          []string{"bash"},
// 		Image:        "all-in-one-ubuntu:latest",
// 	})

// 	service := &Service{manager}
// 	info := CodeExecInfo{
// 		Lang: "Golang",
// 		Content: `package main
// 		import "fmt"

// 		func main() {
// 			fmt.Println("hello world")
// 		}`,
// 		Args: []string{"--debug", "life"},
// 	}
// 	log.Println(info.Content)

// 	res, err := service.executeCode(ctx, info)
// 	if err != nil {
// 		panic(err)
// 	}

// 	log.Println(res)
// }
