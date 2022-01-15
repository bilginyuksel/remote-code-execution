package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codigician/remote-code-execution/internal/handler"
	"github.com/codigician/remote-code-execution/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestBalancer(t *testing.T) {
	e := echo.New()
	mockBalancer := mocks.NewMockBalancerService(gomock.NewController(t))
	balancerHandler := handler.NewBalancer(mockBalancer)
	balancerHandler.RegisterRoutes(e)

	srv := httptest.NewServer(e.Server.Handler)
	defer srv.Close()

	for _, tc := range rceTests {
		t.Run(tc.scenario, func(t *testing.T) {
			reqBytes, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/codexec", srv.URL), bytes.NewBuffer(reqBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(echo.HeaderAccept, echo.MIMETextPlain)

			mockBalancer.EXPECT().
				Exec(gomock.Any(), tc.mockCodexecReq).
				Return([]byte("some response"), tc.mockErr).
				AnyTimes()

			res, _ := http.DefaultClient.Do(req)
			res.Body.Close()

			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)
		})
	}
}
