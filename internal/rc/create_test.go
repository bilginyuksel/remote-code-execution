package rc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/codigician/remote-code-execution/internal/rc"
	"github.com/docker/docker/api/types/container"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreate_ContainerCreateFailed_ReturnErr(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)
	mockContainerPort.EXPECT().
		ContainerCreate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(container.ContainerCreateCreatedBody{}, errors.New("some kind of error"))

	cli := rc.NewClient(mockContainerPort, nil)
	_, err := cli.Create(context.TODO(), &container.HostConfig{})

	assert.NotNil(t, err)
}

func TestCreate_ContainerStartFailed_ReturnErr(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)
	mockContainerPort.EXPECT().
		ContainerCreate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(container.ContainerCreateCreatedBody{ID: "some-id"}, nil)
	mockContainerPort.EXPECT().
		ContainerStart(gomock.Any(), "some-id", gomock.Any()).
		Return(errors.New("err"))

	cli := rc.NewClient(mockContainerPort, nil)
	_, err := cli.Create(context.TODO(), &container.HostConfig{})

	assert.NotNil(t, err)
}

func TestCreate_ContainerCreatedAndStarted_ReturnContainerID(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)
	mockContainerPort.EXPECT().
		ContainerCreate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(container.ContainerCreateCreatedBody{ID: "some-id"}, nil)
	mockContainerPort.EXPECT().
		ContainerStart(gomock.Any(), "some-id", gomock.Any()).
		Return(nil)

	cli := rc.NewClient(mockContainerPort, nil)
	containerID, err := cli.Create(context.TODO(), &container.HostConfig{})

	assert.Nil(t, err)
	assert.Equal(t, "some-id", containerID)
}
