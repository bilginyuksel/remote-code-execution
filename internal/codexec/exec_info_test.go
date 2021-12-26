package codexec_test

import (
	"testing"

	"github.com/codigician/remote-code-execution/internal/codexec"
	"github.com/stretchr/testify/assert"
)

func TestCommand_NotSupportedLanguage_ReturnEmptyCommand(t *testing.T) {
	info := codexec.ExecutionInfo{Lang: "not-exists"}

	assert.Equal(t, []string{""}, info.Command())
}

func TestCommand_SupportedLanguages(t *testing.T) {
	testcases := []struct {
		scenario    string
		info        codexec.ExecutionInfo
		expectedCmd []string
	}{
		{
			scenario:    "Given python3 language",
			info:        codexec.ExecutionInfo{Lang: "python3", Args: []string{"hello", "5", "3"}},
			expectedCmd: []string{"python3", "Main.py", "hello", "5", "3"},
		},
		{
			scenario:    "Given golang language",
			info:        codexec.ExecutionInfo{Lang: "golang"},
			expectedCmd: []string{"/usr/local/go/bin/go", "run", "Main.go"},
		},
		{
			scenario:    "Given java8 language",
			info:        codexec.ExecutionInfo{Lang: "java8"},
			expectedCmd: []string{"javac", "Main.java", "&&", "java", "Main"},
		},
		{
			scenario:    "Given javascript language",
			info:        codexec.ExecutionInfo{Lang: "javascript"},
			expectedCmd: []string{"nodejs", "Main.js"},
		},
		{
			scenario:    "Given nodejs language",
			info:        codexec.ExecutionInfo{Lang: "nodejs"},
			expectedCmd: []string{"nodejs", "Main.js"},
		},
		{
			scenario:    "Given c++ language",
			info:        codexec.ExecutionInfo{Lang: "c++"},
			expectedCmd: []string{"g++", "Main.cpp", "-o", "Main", "&&", "./Main"},
		},
		{
			scenario:    "Given c language",
			info:        codexec.ExecutionInfo{Lang: "c"},
			expectedCmd: []string{"gcc", "Main.c", "-o", "Main", "&&", "./Main"},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.scenario, func(t *testing.T) {
			assert.Equal(t, testcase.expectedCmd, testcase.info.Command())
		})
	}
}
