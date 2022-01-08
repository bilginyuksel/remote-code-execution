package codexec

import (
	"fmt"
	"strings"
)

const _spaceCharacter = " "

type ExecutionInfo struct {
	Lang    string
	Content string
	Args    []string
}

func (info ExecutionInfo) Command() []string {
	if !supportedLanguages.IsSupported(info.Lang) {
		return nil
	}
	characteristic := supportedLanguages.Get(info.Lang)
	// use bash -c to execute multiple commands at once
	// for example: javac Main.java && java Main
	args := strings.Join(info.Args, _spaceCharacter)
	executionCommand := fmt.Sprintf("%s %s", characteristic.Commands, args)
	return []string{"bash", "-c", strings.TrimSpace(executionCommand)}
}
