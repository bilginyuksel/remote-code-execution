package app

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:   "Run Remote Code Execution Engine",
		Usage:  "Run RCEE HTTP application",
		Action: runApp,
	}
}

func runApp(c *cli.Context) error {
	fmt.Println("running the application on port ...")
	return nil
}
