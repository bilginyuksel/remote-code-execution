package codexec

import (
	"strings"
)

type ExecutionInfo struct {
	Lang    string
	Content string
	Args    []string
}

func (info ExecutionInfo) Command() []string {
	characteristic := supportedLanguages.Get(info.Lang)
	cmds := strings.Split(characteristic.Commands, " ")
	return append(cmds, info.Args...)
}
