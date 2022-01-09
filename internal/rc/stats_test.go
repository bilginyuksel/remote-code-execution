package rc_test

import (
	"context"
	"testing"

	"github.com/codigician/remote-code-execution/internal/rc"
	"github.com/golang/mock/gomock"
)

func TestStats_CallContainerStatsOneShot(t *testing.T) {
	mockContainerPort := newMockContainerPort(t)

	mockContainerPort.EXPECT().ContainerStatsOneShot(gomock.Any(), "cid-1")

	client := rc.NewClient(mockContainerPort, nil)
	client.Stats(context.TODO(), "cid-1")
}
