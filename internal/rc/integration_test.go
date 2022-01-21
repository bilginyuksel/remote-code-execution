package rc_test

import (
	"context"
	"log"
	"testing"

	"github.com/codigician/remote-code-execution/internal/rc"
	"github.com/codigician/remote-code-execution/pkg/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/suite"
)

type DockerContainerClientTestSuite struct {
	suite.Suite
	dockerClient        *client.Client
	client              *rc.Client
	createdContainerIDs []string
}

func (s *DockerContainerClientTestSuite) SetupSuite() {
	var containerConfig container.Config
	if err := config.Read("../../.config/test.yml", &containerConfig); err != nil {
		s.Fail(err.Error())
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		s.Fail(err.Error())
	}

	s.dockerClient = dockerClient
	s.client = rc.NewClient(dockerClient, &containerConfig)
}

func (s *DockerContainerClientTestSuite) TearDownSuite() {
	ctx := context.Background()
	for _, id := range s.createdContainerIDs {
		_ = s.dockerClient.ContainerStop(ctx, id, nil)
		_ = s.dockerClient.ContainerRemove(ctx, id, types.ContainerRemoveOptions{})
	}
	s.dockerClient.Close()
}

func TestIntegrationRemoteContainerWithDockerClient(t *testing.T) {
	suite.Run(t, new(DockerContainerClientTestSuite))
}

func (s *DockerContainerClientTestSuite) TestList() {
	l, _ := s.client.List(context.Background())
	log.Println(l)
}

func (s *DockerContainerClientTestSuite) TestCreateContainer() {
	ctx := context.Background()
	cid, err := s.client.Create(ctx, &container.HostConfig{})
	s.createdContainerIDs = append(s.createdContainerIDs, cid)

	clist, _ := s.dockerClient.ContainerList(ctx, types.ContainerListOptions{})

	s.Nil(err)
	s.Equal(cid, clist[0].ID)
}
