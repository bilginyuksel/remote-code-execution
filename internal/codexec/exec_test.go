package codexec_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/codigician/remote-code-execution/internal/codexec"
	"github.com/codigician/remote-code-execution/internal/mocks"
	"github.com/docker/docker/api/types/container"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestExec_NotSupportedLanguage_ReturnErr(t *testing.T) {
	service := codexec.New(newMockContainerClient(t), &container.HostConfig{}, nil)
	_, err := service.Exec(context.TODO(), codexec.ExecutionInfo{Lang: "nolang"})

	assert.Equal(t, "language does not supported", err.Error())
}

func TestExec_ContainerCreateFailure_ReturnErr(t *testing.T) {
	mockContainerClient := newMockContainerClient(t)
	service := codexec.New(mockContainerClient, &container.HostConfig{}, nil)
	mockContainerClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return("", errors.New("create failed"))

	_, err := service.Exec(context.TODO(), codexec.ExecutionInfo{Lang: "python3"})

	assert.Equal(t, "create failed", err.Error())
}

func TestExec_WriteFailure_ReturnErr(t *testing.T) {
	mockContainerClient := newMockContainerClient(t)
	mockContainerClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return("", nil)
	mockContainerClient.EXPECT().ForceRemove(gomock.Any(), gomock.Any())

	service := codexec.New(mockContainerClient, &container.HostConfig{},
		func(baseDir, filename, content string) (string, error) { return "", errors.New("could not write") })

	_, err := service.Exec(context.TODO(), codexec.ExecutionInfo{Lang: "golang"})

	assert.Equal(t, "could not write", err.Error())
}

func TestExec(t *testing.T) {
	containerID := "c1"
	mockResponse := []byte("resp")
	mockFilepath := fmt.Sprintf("%s/%s/Main.go", codexec.MountSource, "ransomid")
	expectedFileDir := fmt.Sprintf("%s/%s/", codexec.MountTarget, "ransomid")
	expectedCmd := []string{"/usr/local/go/bin/go", "run", "Main.go", "yuksel"}

	mockContainerClient := newMockContainerClient(t)
	mockContainerClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(containerID, nil)
	mockContainerClient.EXPECT().ForceRemove(gomock.Any(), containerID)
	mockContainerClient.EXPECT().Exec(gomock.Any(), containerID, expectedFileDir, expectedCmd).Return(mockResponse, nil)

	service := codexec.New(mockContainerClient, &container.HostConfig{},
		func(baseDir, filename, content string) (string, error) { return mockFilepath, nil })

	res, err := service.Exec(context.TODO(), codexec.ExecutionInfo{
		Lang: "golang",
		Content: `package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("greetings to %s\n", os.Args[1])
}`,
		Args: []string{"yuksel"},
	})

	assert.Nil(t, err)
	assert.Equal(t, mockResponse, res)
}

func newMockContainerClient(t *testing.T) *mocks.MockContainerClient {
	return mocks.NewMockContainerClient(gomock.NewController(t))
}