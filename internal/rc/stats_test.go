package rc_test

import (
	"context"
	"testing"

	"github.com/codigician/remote-code-execution/internal/rc"
	"github.com/golang/mock/gomock"
)

func TestStats_CallContainerStatsOneShot(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)
	const containerID = "cid-1"

	mockContainerPort.EXPECT().ContainerStatsOneShot(gomock.Any(), "cid-1")

	client := rc.NewClient(mockContainerPort, nil)
	_, _ = client.Stats(context.TODO(), containerID)
}
