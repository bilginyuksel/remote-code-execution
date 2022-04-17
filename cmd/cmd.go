package cmd

import (
	"github.com/urfave/cli/v2"
)

func Build() *cli.App {
	return &cli.App{
		Name:  "Remote Code Execution",
		Usage: "Remote code execution CLI menu",
		Commands: []*cli.Command{
			CommandExec(),
			CommandServe(),
			CommandGrpcServer(),
		},
	}
}
