package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"

	"github.com/codigician/remote-code-execution/internal/codexec"
	pb "github.com/codigician/remote-code-execution/internal/grpc"
	"github.com/codigician/remote-code-execution/internal/handler"
	"github.com/codigician/remote-code-execution/internal/rc"
	"github.com/codigician/remote-code-execution/pkg/config"
)

type rceService interface {
	Exec(ctx context.Context, info codexec.ExecutionInfo) (*codexec.ExecutionRes, error)
}

type server struct {
	service rceService
	pb.UnimplementedCodeExecutorServiceServer
}

func (s *server) Exec(ctx context.Context, info *pb.CodeExecutionRequest) (*pb.CodeExecutionResponse, error) {
	res, err := s.service.Exec(ctx, codexec.ExecutionInfo{
		Lang:    info.Lang,
		Content: info.Content,
		Args:    strings.Split(info.Args, ","),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CodeExecutionResponse{
		Output:          res.Output,
		ExecutionTimeMs: res.ExecutionTime.Milliseconds(),
	}, nil
}

func CommandGrpcServer() *cli.Command {
	return &cli.Command{
		Name:  "grpc-server",
		Usage: "Run grpc server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   "8080",
				Usage:   "Port to listen on",
			},
		},
		Action: func(c *cli.Context) error {
			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Int("port")))
			if err != nil {
				return err
			}

			var (
				containerConfig     container.Config
				containerHostConfig container.HostConfig
				env                 = os.Getenv("APP_ENV")
			)

			if err := config.Read(fmt.Sprintf(".config/%s.yml", env), &containerConfig, &containerHostConfig); err != nil {
				return err
			}

			dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				return err
			}

			e := echo.New()
			containerClient := rc.NewClient(dockerClient, &containerConfig)

			codexecService := codexec.New(containerClient, &containerHostConfig, codexec.WriteFile)
			codexecHandler := handler.NewRemoteCodeExecutor(codexecService)
			codexecHandler.RegisterRoutes(e)

			pool := codexec.NewContainerPool()
			ticker := time.NewTicker(_balancerIntervalDuration)
			balancerService := codexec.NewContainerBalancer(containerClient, pool, &containerHostConfig, codexecService)
			balancerService.FillPool(context.Background())
			go balancerService.Balance(ticker)

			s := grpc.NewServer()
			pb.RegisterCodeExecutorServiceServer(s, &server{service: balancerService})
			log.Printf("Listening on port %d\n", c.Int("port"))
			return s.Serve(lis)
		},
	}
}
