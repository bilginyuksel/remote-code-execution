package main

import (
	"context"
	"log"
)

type BalancerConfiguration struct {
	MaxContainerCount int
	MinContainerCount int
}

type BalancerContainerManager interface {
	Create(ctx context.Context) (string, error)
	//	Status(ctx context.Context, containerID string) (ContainerStatus, error)
	Remove(ctx context.Context)
}

type Balancer struct {
	conf               BalancerConfiguration
	manager            BalancerContainerManager
	currContainerCount int
}

func (b *Balancer) Balance(ctx context.Context) (createdContainers []string, err error) {
	if b.currContainerCount > b.conf.MinContainerCount {
		log.Printf("no need to rebalance there are already %d containers\n", b.currContainerCount)
		return nil, nil
	}

	for ; b.currContainerCount < b.conf.MinContainerCount; b.currContainerCount++ {
		containerID, err := b.manager.Create(ctx)
		if err != nil {
			log.Printf("create container: %v\n", err)
			return createdContainers, err
		}
		createdContainers = append(createdContainers, containerID)
	}

	return createdContainers, err
}
