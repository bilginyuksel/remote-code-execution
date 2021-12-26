package main

import (
	"os"

	"github.com/codigician/remote-code-execution/cmd"
)

func main() {
	app := cmd.Build()
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
