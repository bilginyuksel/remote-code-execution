package rc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/codigician/remote-code-execution/internal/rc"
	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
)

func TestForceRemove(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)
	mockContainerPort.EXPECT().ContainerStop(gomock.Any(), "container-to-remove", nil).Return(nil)
	mockContainerPort.EXPECT().ContainerRemove(gomock.Any(), "container-to-remove", types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}).Return(nil)

	client := rc.NewClient(mockContainerPort, nil)
	client.ForceRemove(context.TODO(), "container-to-remove")
}

func TestForceRemove_Failed_NothingHappens(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)
	mockContainerPort.EXPECT().ContainerStop(gomock.Any(), "container-to-remove", nil).Return(nil)
	mockContainerPort.EXPECT().
		ContainerRemove(gomock.Any(), "container-to-remove", gomock.Any()).
		Return(errors.New("some err"))

	client := rc.NewClient(mockContainerPort, nil)
	client.ForceRemove(context.TODO(), "container-to-remove")
}

func TestForceRemove_ContainerStopFailed_CallContainerRemoveAnyways(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)
	mockContainerPort.EXPECT().ContainerStop(gomock.Any(), "c1", nil).Return(errors.New("some err"))
	mockContainerPort.EXPECT().ContainerRemove(gomock.Any(), "c1", gomock.Any()).Return(nil)

	client := rc.NewClient(mockContainerPort, nil)
	client.ForceRemove(context.TODO(), "c1")
}
