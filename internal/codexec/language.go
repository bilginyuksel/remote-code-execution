package codexec

import (
	"fmt"
	"strings"
)

const ProgramFilename = "Main"

var (
	Golang     = newLanguageCharacteristic("golang", "/usr/local/go/bin/go run Main.go", ".go")
	Python3    = newLanguageCharacteristic("python3", "python3 Main.py", ".py")
	Javascript = newLanguageCharacteristic("javascript", "nodejs Main.js", ".js")
	NodeJS     = newLanguageCharacteristic("nodejs", "nodejs Main.js", ".js")
	Java8      = newLanguageCharacteristic("java8", "javac Main.java && java Main", ".java")
	C          = newLanguageCharacteristic("c", "gcc Main.c -o Main && ./Main", ".c")
	CPlusplus  = newLanguageCharacteristic("c++", "g++ Main.cpp -o Main && ./Main", ".cpp")

	supportedLanguages = SupportedLanguages{
		"golang":     Golang,     // /usr/local/go/bin/go <filename> <arg1> <arg2>
		"python3":    Python3,    // python3 <filename> <arg1> <arg2>
		"javascript": Javascript, // nodejs <filename> <arg1> <arg2>
		"nodejs":     NodeJS,     // nodejs <filename> <arg1> <arg2>
		"java8":      Java8,      // javac <filename> && java <filename> <arg1> <arg2>
		"c":          C,          // gcc <filename> && <executable> <arg1> <arg2>
		"c++":        CPlusplus,  // gcc <filename> && <executable> <arg1> <arg2>
	}
)

type (
	SupportedLanguages map[string]LanguageCharacteristic

	LanguageCharacteristic struct {
		Name      string
		Commands  string
		Extension string
	}
)

func (s SupportedLanguages) IsSupported(lang string) bool {
	_, ok := s[strings.ToLower(lang)]
	return ok
}

func (s SupportedLanguages) Get(lang string) LanguageCharacteristic {
	return s[strings.ToLower(lang)]
}

func (c LanguageCharacteristic) Filename() string {
	return fmt.Sprintf("%s%s", ProgramFilename, c.Extension)
}

func newLanguageCharacteristic(name, cmds, extension string) LanguageCharacteristic {
	return LanguageCharacteristic{
		Name:      name,
		Commands:  cmds,
		Extension: extension,
	}
}
