package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codigician/remote-code-execution/internal/codexec"
	"github.com/codigician/remote-code-execution/internal/handler"
	"github.com/codigician/remote-code-execution/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var rceTests = []struct {
	scenario           string
	requestBody        interface{}
	mockCodexecReq     codexec.ExecutionInfo
	mockErr            error
	expectedStatusCode int
}{
	{
		scenario:           "Given unexpected request type return bad request",
		requestBody:        "unexpected request",
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		scenario:           "Given a request, when the service crashed return internal server error",
		requestBody:        handler.RemoteCodeExecuteRequest{Lang: "python3", Content: "print('hello world')"},
		expectedStatusCode: http.StatusInternalServerError,
		mockErr:            errors.New("some kind of error happened"),
		mockCodexecReq:     codexec.ExecutionInfo{Lang: "python3", Content: "print('hello world')"},
	},
	{
		scenario:           "Given a request, execute the code and return the response from service",
		requestBody:        handler.RemoteCodeExecuteRequest{Lang: "nodejs", Content: "console.log('hello world')"},
		expectedStatusCode: http.StatusOK,
		mockCodexecReq:     codexec.ExecutionInfo{Lang: "nodejs", Content: "console.log('hello world')"},
	},
}

func TestRemoteCodeExecutor(t *testing.T) {
	e := echo.New()
	mockExecutor := newMockExecutor(t)
	rceHandler := handler.NewRemoteCodeExecutor(mockExecutor)
	rceHandler.RegisterRoutes(e)

	srv := httptest.NewServer(e.Server.Handler)
	defer srv.Close()

	for _, tc := range rceTests {
		t.Run(tc.scenario, func(t *testing.T) {
			reqBytes, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/codexec", srv.URL), bytes.NewBuffer(reqBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(echo.HeaderAccept, echo.MIMETextPlain)

			mockExecutor.EXPECT().
				ExecOnce(gomock.Any(), tc.mockCodexecReq).
				Return([]byte("some response"), tc.mockErr).
				AnyTimes()

			res, _ := http.DefaultClient.Do(req)
			res.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)
		})
	}
}

func newMockExecutor(t *testing.T) *mocks.MockRemoteCodeExecutorService {
	return mocks.NewMockRemoteCodeExecutorService(gomock.NewController(t))
}
