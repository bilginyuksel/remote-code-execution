package handler

import (
	"context"
	"net/http"

	"github.com/codigician/remote-code-execution/internal/codexec"
	"github.com/labstack/echo/v4"
)

type (
	RemoteCodeExecutorService interface {
		ExecOnce(ctx context.Context, info codexec.ExecutionInfo) (*codexec.ExecutionRes, error)
	}

	RemoteCodeExecutor struct {
		rce RemoteCodeExecutorService
	}

	RemoteCodeExecuteRequest struct {
		Lang    string   `json:"lang"`
		Content string   `json:"content"`
		Args    []string `json:"args"`
	}

	RemoteCodeExecutionResponse struct {
		Output                    string `json:"output"`
		ExecutionTimeMilliseconds int64  `json:"execution_time_ms"`
	}
)

func NewRemoteCodeExecutor(rce RemoteCodeExecutorService) *RemoteCodeExecutor {
	return &RemoteCodeExecutor{rce}
}

func (r *RemoteCodeExecutor) RegisterRoutes(e *echo.Echo) {
	e.POST("/v1/codexec", r.Exec)
}

func (r *RemoteCodeExecutor) Exec(ctx echo.Context) error {
	var request RemoteCodeExecuteRequest
	if err := ctx.Bind(&request); err != nil {
		return err
	}

	res, err := r.rce.ExecOnce(ctx.Request().Context(), request.ToCodexecInfo())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, FromExecutionResult(res))
}

func (r RemoteCodeExecuteRequest) ToCodexecInfo() codexec.ExecutionInfo {
	return codexec.ExecutionInfo{
		Lang:    r.Lang,
		Content: r.Content,
		Args:    r.Args,
	}
}

func FromExecutionResult(res *codexec.ExecutionRes) *RemoteCodeExecutionResponse {
	return &RemoteCodeExecutionResponse{
		Output:                    res.Output,
		ExecutionTimeMilliseconds: res.ExecutionTime.Milliseconds(),
	}
}
