package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/stdcopy"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type Container struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

type ContainerStdout []byte

func (c ContainerStdout) Unmarshal(value interface{}) error {
	return json.Unmarshal(c, value)
}

type ContainerPort interface {
	ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error)
	ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error)
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.ContainerCreateCreatedBody, error)
	ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error
	ContainerExecCreate(ctx context.Context, id string, options types.ExecConfig) (types.IDResponse, error)
	ContainerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error)
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
		log.Println("image does not exists")
		if err := c.createImage(ctx, c.defaultImage); err != nil {
			log.Printf("pull image: %v", err)
			return nil, err
		}
	}
	log.Println("given image ready..")
	containerID, err := c.createContainer(ctx)
	if err != nil {
		return nil, err
	}

	if err := c.containerPort.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	return &Container{
		ID:        containerID,
		CreatedAt: time.Now(),
	}, nil
}

func (c *ContainerManager) Exec(ctx context.Context, id string, cmd []string) (ContainerStdout, error) {
	createExecResponse, err := c.containerPort.ContainerExecCreate(ctx, id, types.ExecConfig{
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
	if _, err := stdcopy.StdCopy(&outBuffer, &errBuffer, attachExecRes.Conn); err != nil {
		return nil, err
	}

	return outBuffer.Bytes(), nil
}

func (c *ContainerManager) isImageExists(ctx context.Context, image string) bool {
	imageSummary, err := c.containerPort.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		log.Printf("image list: %v\n", err)
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

	log.Printf("warnings: %v\n", container.Warnings)
	return container.ID, nil
}
