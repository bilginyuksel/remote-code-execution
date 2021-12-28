package config_test

import (
	"testing"

	"github.com/codigician/remote-code-execution/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestRead_BadConfig_ReturnErr(t *testing.T) {
	err := config.Read("path-dont-exists", nil)

	assert.NotNil(t, err)
}

func TestRead_GoodConfig_FillStruct(t *testing.T) {
	type Test struct {
		Appname string
		Host    string
		Port    int
	}
	var actual Test
	err := config.Read("../../testdata/test.yml", &actual)

	assert.Nil(t, err)
	assert.Equal(t, Test{
		Appname: "test-app",
		Host:    "localhost",
		Port:    8000,
	}, actual)
}

func TestRead_InvalidConfigDataToFill_UnmarshalReturnErr(t *testing.T) {
	var invalidConf int

	err := config.Read("../../testdata/test.yml", &invalidConf)

	assert.NotNil(t, err)
}
