package handler

import (
	"context"
	"net/http"

	"github.com/codigician/remote-code-execution/internal/codexec"
	"github.com/labstack/echo/v4"
)

type (
	RemoteCodeExecutorService interface {
		Exec(ctx context.Context, info codexec.ExecutionInfo) ([]byte, error)
	}

	RemoteCodeExecutor struct {
		rce RemoteCodeExecutorService
	}

	RemoteCodeExecuteRequest struct {
		Lang    string   `json:"lang"`
		Content string   `json:"content"`
		Args    []string `json:"args"`
	}
)

func NewRemoteCodeExecutor(rce RemoteCodeExecutorService) *RemoteCodeExecutor {
	return &RemoteCodeExecutor{rce}
}

func (r *RemoteCodeExecutor) RegisterRoutes(e *echo.Echo) {
	e.POST("/codexec", r.Exec)
}

func (r *RemoteCodeExecutor) Exec(ctx echo.Context) error {
	var request RemoteCodeExecuteRequest
	if err := ctx.Bind(&request); err != nil {
		return err
	}

	res, err := r.rce.Exec(ctx.Request().Context(), request.ToCodexecInfo())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Blob(http.StatusOK, echo.MIMETextPlain, res)
}

func (r RemoteCodeExecuteRequest) ToCodexecInfo() codexec.ExecutionInfo {
	return codexec.ExecutionInfo{
		Lang:    r.Lang,
		Content: r.Content,
		Args:    r.Args,
	}
}
