package codexec_test

import (
	"os"
	"testing"

	"github.com/codigician/remote-code-execution/internal/codexec"
	"github.com/stretchr/testify/assert"
)

func TestWriteFile(t *testing.T) {
	_ = os.Mkdir("testdata", 0777)
	filepath, err := codexec.WriteFile("./testdata", "somefile", "")

	assert.Nil(t, err)
	assert.FileExists(t, filepath)

	os.RemoveAll("testdata")
}
