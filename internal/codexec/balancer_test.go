package codexec_test

import (
	"context"
	"errors"
	"testing"

	"github.com/codigician/remote-code-execution/internal/codexec"
	"github.com/codigician/remote-code-execution/internal/mocks"
	"github.com/docker/docker/api/types/container"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFillPool(t *testing.T) {
	mockContainerClient := newMockContainerClient(t)
	pool := codexec.NewContainerPool()
	balancer := codexec.NewContainerBalancer(mockContainerClient, pool, nil, nil)

	mockContainerClient.EXPECT().
		List(gomock.Any()).
		Return([]string{"cid-1", "cid-2"}, nil)

	balancer.FillPool(context.Background())

	assert.Contains(t, pool.Nodes, "cid-1")
	assert.Contains(t, pool.Nodes, "cid-2")
}

func TestFillPool_ContainerClientListErr_DontFillPool(t *testing.T) {
	mockContainerClient := newMockContainerClient(t)
	pool := codexec.NewContainerPool()
	balancer := codexec.NewContainerBalancer(mockContainerClient, pool, nil, nil)

	mockContainerClient.EXPECT().
		List(gomock.Any()).
		Return(nil, errors.New("list error"))

	balancer.FillPool(context.Background())

	assert.Len(t, pool.Nodes, 0)
}

func TestBalancerExec_NoContainersInPool_CreateContainerAddToPoolThenExecute(t *testing.T) {
	mockContainerClient := newMockContainerClient(t)
	mockCodexecutor := mocks.NewMockCodexecutor(gomock.NewController(t))
	pool := codexec.NewContainerPool()
	conf := &container.HostConfig{}
	balancer := codexec.NewContainerBalancer(mockContainerClient, pool, conf, mockCodexecutor)

	info := codexec.ExecutionInfo{}

	mockContainerClient.EXPECT().Create(context.Background(), conf).Return("cid-1", nil).Times(1)
	mockCodexecutor.EXPECT().Exec(gomock.Any(), "cid-1", info).Return([]byte("result"), nil).Times(1)

	res, err := balancer.Exec(context.Background(), info)

	assert.NotNil(t, res)
	assert.Nil(t, err)
	assert.Equal(t, "cid-1", pool.Get().ID)
}

func TestBalancerExec_NoContainersInPoolCreateContainerFails_ReturnErr(t *testing.T) {
	mockContainerClient := newMockContainerClient(t)
	pool := codexec.NewContainerPool()

	info := codexec.ExecutionInfo{Lang: "python3"}
	hostConfig := &container.HostConfig{}

	mockContainerClient.EXPECT().Create(gomock.Any(), hostConfig).
		Return("", errors.New("some error happened"))

	balancer := codexec.NewContainerBalancer(mockContainerClient, pool, hostConfig, nil)
	res, err := balancer.Exec(context.Background(), info)

	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestBalancerExec_ThereIsContainerInPool_GetContainerThenExecute(t *testing.T) {
	mockCodexecutor := mocks.NewMockCodexecutor(gomock.NewController(t))
	pool := codexec.NewContainerPool()
	pool.Add("cid-1")

	var info codexec.ExecutionInfo
	mockCodexecutor.EXPECT().Exec(gomock.Any(), "cid-1", info).Return([]byte("result"), nil).Times(1)

	balancer := codexec.NewContainerBalancer(nil, pool, nil, mockCodexecutor)
	res, err := balancer.Exec(context.Background(), info)

	assert.Nil(t, err)
	assert.NotNil(t, res)
}
