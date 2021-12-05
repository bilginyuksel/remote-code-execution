package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type Container struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

type ContainerPort interface {
	ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error)
	ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error)
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.ContainerCreateCreatedBody, error)
	ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error
}

type ContainerManager struct {
	defaultImage  string
	containerPort ContainerPort
}

func NewManager(containerPort ContainerPort, image string) *ContainerManager {
	return &ContainerManager{
		defaultImage:  image,
		containerPort: containerPort,
	}
}

func (c *ContainerManager) Create(ctx context.Context) (*Container, error) {
	if !c.isImageExists(ctx, c.defaultImage) {
		if err := c.createImage(ctx, c.defaultImage); err != nil {
			return nil, err
		}
	}
	containerID, err := c.createContainer(ctx)
	if err != nil {
		return nil, err
	}

	if err = c.startContainer(ctx, containerID); err != nil {
		return nil, err
	}

	return &Container{
		ID:        containerID,
		CreatedAt: time.Now(),
	}, nil
}

func (c *ContainerManager) isImageExists(ctx context.Context, image string) bool {
	imageSummary, err := c.containerPort.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return false
	}

	for idx := range imageSummary {
		if _, ok := imageSummary[idx].Labels[image]; ok {
			return true
		}
	}

	return false
}

func (c *ContainerManager) createImage(ctx context.Context, image string) error {
	reader, err := c.containerPort.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	readerBytes, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	log.Println(string(readerBytes))

	return err
}

func (c *ContainerManager) createContainer(ctx context.Context) (string, error) {
	container, err := c.containerPort.ContainerCreate(ctx, &container.Config{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{"bash"},
		Image:        "ubuntu",
	}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	log.Println(container.Warnings)
	return container.ID, nil
}

func (c *ContainerManager) startContainer(ctx context.Context, id string) error {
	return c.containerPort.ContainerStart(ctx, id, types.ContainerStartOptions{})
}
