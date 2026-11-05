package handler

import (
	"net/http"
	"testing"

	"github.com/datarhei/core/v16/http/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func getDummyPingRouter() *echo.Echo {
	router := mock.DummyEcho()

	handler := NewPing()

	router.Add("GET", "/", handler.Ping)

	return router
}

func TestPing(t *testing.T) {
	router := getDummyPingRouter()

	response := mock.Request(t, http.StatusOK, router, "GET", "/", nil)

	require.Equal(t, "pong", string(response.Data.([]byte)))
}
