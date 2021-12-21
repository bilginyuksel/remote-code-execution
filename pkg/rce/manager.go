package rce

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/stdcopy"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type ContainerPort interface {
	ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error)
	ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error)
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.ContainerCreateCreatedBody, error)
	ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error
	ContainerExecCreate(ctx context.Context, id string, options types.ExecConfig) (types.IDResponse, error)
	ContainerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error)
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
}

type ContainerManager struct {
	containerConfig *container.Config
	containerPort   ContainerPort
}

// NewContainerManager create the conatiner manager with the given container configurations
func NewContainerManager(containerPort ContainerPort, containerConfig *container.Config) *ContainerManager {
	return &ContainerManager{
		containerPort:   containerPort,
		containerConfig: containerConfig,
	}
}

// Create creates a container then returns the container id
func (c *ContainerManager) Create(ctx context.Context) (string, error) {
	if err := c.CreateImageIfNotExists(ctx, c.containerConfig.Image); err != nil {
		return "", err
	}

	container, err := c.containerPort.ContainerCreate(ctx, c.containerConfig, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:     mount.TypeBind,
				Source:   "/Users/bilginyuksel/ps-workspace/remote-code-execution/example",
				Target:   "/app",
				ReadOnly: true,
			},
		},
	}, nil, nil, "")
	if err != nil {
		return "", err
	}
	log.Printf("warnings: %v\n", container.Warnings)

	if err := c.containerPort.ContainerStart(ctx, container.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	return container.ID, nil
}

// List list the containers
func (c *ContainerManager) List(ctx context.Context) ([]string, error) {
	containers, err := c.containerPort.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	var containerIDs []string
	for idx := range containers {
		containerIDs = append(containerIDs, containers[idx].ID)
	}

	return containerIDs, nil
}

// Exec executes the given command to the container with the given id
func (c *ContainerManager) Exec(ctx context.Context, containerID string, cmd []string) ([]byte, error) {
	createExecResponse, err := c.containerPort.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		Cmd:          cmd,
		AttachStderr: true,
		AttachStdout: true,
	})
	if err != nil {
		log.Printf("container exec create: %v, cmd: %v\n", err, cmd)
		return nil, err
	}
	log.Println("container exec created")

	attachExecRes, err := c.containerPort.ContainerExecAttach(ctx, createExecResponse.ID, types.ExecStartCheck{})
	if err != nil {
		log.Printf("container exec attach: %v\n", err)
		return nil, err
	}
	log.Println("container exec attached")

	var outBuffer, errBuffer bytes.Buffer
	_, err = stdcopy.StdCopy(&outBuffer, &errBuffer, attachExecRes.Conn)

	return outBuffer.Bytes(), err
}

func (c *ContainerManager) CreateImageIfNotExists(ctx context.Context, image string) error {
	if c.IsImageExists(ctx, c.containerConfig.Image) {
		log.Println("image already exists")
		return nil
	}

	return c.PullImage(ctx, image)
}

// IsImageExists get the image list and compare the images with the given image
// if the image is in the list return true otherwise return false
func (c *ContainerManager) IsImageExists(ctx context.Context, image string) bool {
	imageSummary, err := c.containerPort.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		log.Printf("image list failed: %v\n", err)
		return false
	}

	for idx := range imageSummary {
		repoTags := imageSummary[idx].RepoTags
		for _, tag := range repoTags {
			if tag == image {
				return true
			}
		}
	}

	return false
}

func (c *ContainerManager) PullImage(ctx context.Context, image string) error {
	log.Println("pulling image")
	reader, err := c.containerPort.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		log.Printf("failed while pulling the image, err: %v\n", err)
		return err
	}
	defer reader.Close()

	readerBytes, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("failed when reading the bytes, err: %v\n", err)
		return err
	}
	log.Printf("image pull response: %v\n", string(readerBytes))

	return err
}
