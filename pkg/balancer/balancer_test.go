package main

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestBalancer(t *testing.T) {
	mock := NewMockBalancerContainerManager(gomock.NewController(t))

	b := Balancer{
		conf: BalancerConfiguration{
			MaxContainerCount: 5,
			MinContainerCount: 3,
		},
		manager: mock,
	}

	mock.EXPECT().Create(gomock.Any()).Times(3)

	b.Balance(context.Background())
}
