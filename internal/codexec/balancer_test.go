package codexec_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

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

func TestBalancerShutdown(t *testing.T) {
	mockContainerClient := newMockContainerClient(t)
	pool := codexec.NewContainerPool()
	pool.Add("cid-1")
	pool.Add("cid-2")
	pool.Add("cid-3")

	mockContainerClient.EXPECT().ForceRemove(context.TODO(), "cid-1").Times(1)
	mockContainerClient.EXPECT().ForceRemove(context.TODO(), "cid-2").Times(1)
	mockContainerClient.EXPECT().ForceRemove(context.TODO(), "cid-3").Times(1)

	balancer := codexec.NewContainerBalancer(mockContainerClient, pool, nil, nil)
	balancer.Shutdown(context.TODO())

	assert.Empty(t, pool.Nodes)
	assert.Nil(t, pool.Head)
	assert.Nil(t, pool.Tail)
	assert.Nil(t, pool.Curr)
}

func TestBalancerBalance_FailNTimes_CreateNotFailedContainersOnly(t *testing.T) {
	mockContainerClient := newMockContainerClient(t)
	pool := codexec.NewContainerPool()
	balancer := codexec.NewContainerBalancer(mockContainerClient, pool, nil, nil)

	var (
		mockContainerIDs  []string
		mockContainerErrs []error
		callIndex         int
	)
	// 7 successful calls
	for i := 0; i < 7; i++ {
		mockContainerIDs = append(mockContainerIDs, fmt.Sprintf("cid-%d", i))
		mockContainerErrs = append(mockContainerErrs, nil)
	}
	// 3 failed calls
	for i := 0; i < 3; i++ {
		mockContainerIDs = append(mockContainerIDs, "")
		mockContainerErrs = append(mockContainerErrs, errors.New("err"))
	}

	mockContainerClient.EXPECT().Create(context.TODO(), nil).
		DoAndReturn(func(ctx context.Context, conf *container.HostConfig) (string, error) {
			defer func() { callIndex++ }()
			return mockContainerIDs[callIndex], mockContainerErrs[callIndex]
		}).Times(10)

	ticker := time.NewTicker(100 * time.Millisecond)
	go balancer.Balance(ticker)
	time.Sleep(180 * time.Millisecond)
	ticker.Stop()

	assert.Len(t, pool.Nodes, 7)
}
