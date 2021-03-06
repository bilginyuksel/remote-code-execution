package codexec

import (
	"context"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
)

const _desiredActiveContainerCount = 10

type (
	Codexecutor interface {
		Exec(ctx context.Context, containerID string, info ExecutionInfo) (*ExecutionRes, error)
	}

	ContainerBalancer struct {
		containerPool   *ContainerPool
		hostConfig      *container.HostConfig
		containerClient ContainerClient
		executor        Codexecutor
	}
)

func NewContainerBalancer(containerClient ContainerClient, containerPool *ContainerPool,
	hostConfig *container.HostConfig, executor Codexecutor) *ContainerBalancer {
	return &ContainerBalancer{
		containerPool:   containerPool,
		containerClient: containerClient,
		executor:        executor,
		hostConfig:      hostConfig,
	}
}

// Exec get the execution info and find the next container
// to use that container to execute the programming content
// when execution is done send response to client and start a goroutine
// to update the container status and metrics
func (cb *ContainerBalancer) Exec(ctx context.Context, info ExecutionInfo) (*ExecutionRes, error) {
	if node := cb.containerPool.Get(); node != nil {
		return cb.executor.Exec(ctx, node.ID, info)
	}

	containerID, err := cb.containerClient.Create(ctx, cb.hostConfig)
	if err != nil {
		return nil, err
	}
	cb.containerPool.Add(containerID)
	node := cb.containerPool.Get()
	return cb.executor.Exec(ctx, node.ID, info)
}

// FillPool use container client to list all active containers
// and the fill the pool with the active containers
func (cb *ContainerBalancer) FillPool(ctx context.Context) {
	containerIDs, err := cb.containerClient.List(ctx)
	if err != nil {
		log.Printf("container client list failed, err: %v\n", err)
		return
	}

	for _, containerID := range containerIDs {
		cb.containerPool.Add(containerID)
	}
}

// Balance get the snapshot of the containers average, min and max response times.
// According to containers response times decide to create new containers or delete
// the existing ones to not waste a resource.
// If there is no snapshot just took a snapshot and also
// limit the container creation and deletion to the min and max limits.
func (cb *ContainerBalancer) Balance(ticker *time.Ticker) {
	for range ticker.C {
		log.Println("Rebalancing containers..")
		countOfActiveContainers := len(cb.containerPool.Nodes)
		countOfContainerToCreate := _desiredActiveContainerCount - countOfActiveContainers
		for ; countOfContainerToCreate > 0; countOfContainerToCreate-- {
			id, err := cb.containerClient.Create(context.Background(), cb.hostConfig)
			if err != nil {
				log.Printf("could not create container, err: %v\n", err)
				continue
			}
			cb.containerPool.Add(id)
		}
	}
}

func (cb *ContainerBalancer) Shutdown(ctx context.Context) {
	for node := cb.containerPool.Get(); node != nil; node = cb.containerPool.Get() {
		log.Printf("removing container with the id: %v", node.ID)
		cb.containerClient.ForceRemove(ctx, node.ID)
		cb.containerPool.Remove(node.ID)
	}
}
