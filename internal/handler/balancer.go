package handler

import (
	"context"
	"net/http"

	"github.com/codigician/remote-code-execution/internal/codexec"
	"github.com/labstack/echo/v4"
)

type (
	BalancerService interface {
		Exec(ctx context.Context, info codexec.ExecutionInfo) (*codexec.ExecutionRes, error)
	}

	Balancer struct {
		service BalancerService
	}
)

func NewBalancer(service BalancerService) *Balancer {
	return &Balancer{service}
}

func (b *Balancer) RegisterRoutes(e *echo.Echo) {
	e.POST("/v2/codexec", b.Exec)
}

func (b *Balancer) Exec(ctx echo.Context) error {
	var request RemoteCodeExecuteRequest
	if err := ctx.Bind(&request); err != nil {
		return err
	}

	res, err := b.service.Exec(ctx.Request().Context(), request.ToCodexecInfo())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, FromExecutionResult(res))
}
