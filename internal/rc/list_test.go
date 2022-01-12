package rc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/codigician/remote-code-execution/internal/mocks"
	"github.com/codigician/remote-code-execution/internal/rc"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestList_ContainerListFailed_ReturnErr(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)
	mockContainerPort.EXPECT().ContainerList(gomock.Any(), gomock.Any()).Return(nil, errors.New("fail"))

	client := rc.NewClient(mockContainerPort, &container.Config{})
	containers, err := client.List(context.Background())

	assert.Nil(t, containers)
	assert.NotNil(t, err)
}

func TestList_Listed_ReturnContainerIDs(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)
	expectedListOpts := types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.KeyValuePair{Key: "ancestor", Value: "my-image:latest"},
			filters.KeyValuePair{Key: "status", Value: "running"}),
	}
	mockContainerPort.EXPECT().ContainerList(gomock.Any(), expectedListOpts).
		Return([]types.Container{{ID: "container-1", Image: "my-image:latest"},
			{ID: "container-2", Image: "my-image:latest"}}, nil)

	client := rc.NewClient(mockContainerPort, &container.Config{Image: "my-image:latest"})
	containerIDs, _ := client.List(context.Background())

	assert.Equal(t, []string{"container-1", "container-2"}, containerIDs)
}

func newMockContainerPort(t *testing.T) *mocks.MockContainerPort {
	return mocks.NewMockContainerPort(gomock.NewController(t))
}
