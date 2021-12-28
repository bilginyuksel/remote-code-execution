package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/codigician/remote-code-execution/internal/codexec"
	"github.com/codigician/remote-code-execution/internal/rc"
	"github.com/codigician/remote-code-execution/pkg/config"
	"github.com/docker/docker/api/types/container"
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
	var (
		containerConfig     container.Config
		containerHostConfig container.HostConfig

		env      = os.Getenv("APP_ENV")
		filepath = c.String("path")
		lang     = c.String("lang")
	)

	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	if err := config.Read(fmt.Sprintf(".config/%s.yml", env), &containerConfig, &containerHostConfig); err != nil {
		panic(err)
	}

	rcClient := rc.NewClient(dockerClient, &containerConfig)
	service := codexec.New(rcClient, &containerHostConfig, codexec.WriteFile)
	res, err := service.Exec(c.Context, codexec.ExecutionInfo{Lang: lang, Content: string(content)})
	log.Println(string(res))
	return err
}
