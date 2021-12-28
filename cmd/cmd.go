package cmd

import (
	"github.com/urfave/cli/v2"
)

func Build() *cli.App {
	return &cli.App{
		Name:  "rce",
		Usage: "Remote code execution CLI menu",
		Commands: []*cli.Command{
			CommandApplication(),
		},
	}
}
