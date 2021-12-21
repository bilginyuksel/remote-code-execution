package cmd

import (
	"github.com/codigician/remote-code-execution/cmd/app"
	"github.com/urfave/cli/v2"
)

func Build() *cli.App {
	return &cli.App{
		Name:  "Remote Code Execution Application",
		Usage: "Remote code execution CLI menu",
		Commands: []*cli.Command{
			app.NewCommand(),
		},
	}
}
