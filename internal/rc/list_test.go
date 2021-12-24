package rc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/codigician/remote-code-execution/internal/mocks"
	"github.com/codigician/remote-code-execution/internal/rc"
	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestList_ContainerListFailed_ReturnErr(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)
	mockContainerPort.EXPECT().ContainerList(gomock.Any(), gomock.Any()).Return(nil, errors.New("fail"))

	client := rc.NewClient(mockContainerPort, nil)
	containers, err := client.List(context.Background())

	assert.Nil(t, containers)
	assert.NotNil(t, err)
}

func TestList_Listed_ReturnContainerIDs(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)
	mockContainerPort.EXPECT().ContainerList(gomock.Any(), gomock.Any()).
		Return([]types.Container{{ID: "container-1"}, {ID: "container-2"}}, nil)

	client := rc.NewClient(mockContainerPort, nil)
	containerIDs, _ := client.List(context.Background())

	assert.Equal(t, containerIDs, []string{"container-1", "container-2"})
}

func newMockContainerPort(t *testing.T) *mocks.MockContainerPort {
	return mocks.NewMockContainerPort(gomock.NewController(t))
}
