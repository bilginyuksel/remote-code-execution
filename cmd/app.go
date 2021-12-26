package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/codigician/remote-code-execution/internal/codexec"
	"github.com/codigician/remote-code-execution/internal/rc"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/urfave/cli/v2"
)

func CommandApplication() *cli.Command {
	return &cli.Command{
		Name:  "exec",
		Usage: "Run RCEE HTTP application",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "language",
				Aliases:  []string{"lang", "l"},
				Usage:    "Programming language you want to run",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "Filepath of the file that contains the application code",
				Required: true,
			},
		},
		Action: runApp,
	}
}

func runApp(c *cli.Context) error {
	fmt.Println("running the application on port ...")
	filepath := c.String("path")
	lang := c.String("lang")

	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	containerPort := rc.NewClient(dockerClient, &container.Config{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"bash"},
		Image:        "all-in-one-ubuntu:latest",
	})

	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	service := codexec.New(containerPort, &container.HostConfig{
		Mounts: []mount.Mount{{
			Type:     mount.TypeBind,
			Source:   fmt.Sprintf("%s/target", currDir),
			Target:   "/app",
			ReadOnly: false,
		}},
	}, codexec.WriteFile)
	res, err := service.Exec(c.Context, codexec.ExecutionInfo{
		Lang:    lang,
		Content: string(content),
	})
	if err != nil {
		return err
	}

	log.Println(string(res))
	return nil
}
