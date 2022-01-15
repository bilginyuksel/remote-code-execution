package cmd

import (
	"fmt"
	"os"

	"github.com/codigician/remote-code-execution/internal/codexec"
	"github.com/codigician/remote-code-execution/internal/handler"
	"github.com/codigician/remote-code-execution/internal/rc"
	"github.com/codigician/remote-code-execution/pkg/config"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	"github.com/urfave/cli/v2"
)

func CommandServe() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Run RCEE HTTP application",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "Give the port that you want to run the application",
				Value:   8888,
			},
		},
		Action: startServer,
	}
}

func startServer(c *cli.Context) error {
	var (
		containerConfig     container.Config
		containerHostConfig container.HostConfig
		env                 = os.Getenv("APP_ENV")
	)

	if err := config.Read(fmt.Sprintf(".config/%s.yml", env), &containerConfig, &containerHostConfig); err != nil {
		return err
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	e := echo.New()
	containerClient := rc.NewClient(dockerClient, &containerConfig)

	codexecService := codexec.New(containerClient, &containerHostConfig, codexec.WriteFile)
	codexecHandler := handler.NewRemoteCodeExecutor(codexecService)
	codexecHandler.RegisterRoutes(e)

	return e.Start(fmt.Sprintf(":%d", c.Int("port")))
}
