package api

import (
	"net/http"
	"testing"

	"github.com/datarhei/core/v16/http/api"
	"github.com/datarhei/core/v16/http/mock"
	"github.com/datarhei/core/v16/session"

	"github.com/labstack/echo/v4"
)

func getDummySessionRouter() *echo.Echo {
	router := mock.DummyEcho()

	registry, _ := session.New(session.Config{})
	registry.Register("foo", session.CollectorConfig{})

	collector := registry.Collector("foo")
	collector.RegisterAndActivate("foobar", "", "any", "any")

	handler := NewSession(registry)

	router.Add("GET", "/summary", handler.Summary)
	router.Add("GET", "/active", handler.Active)

	return router
}

func TestSessionSummary(t *testing.T) {
	router := getDummySessionRouter()

	response := mock.Request(t, http.StatusOK, router, "GET", "/summary?collectors=foo", nil)

	mock.Validate(t, api.SessionsSummary{}, response.Data)
}

func TestSessionActive(t *testing.T) {
	router := getDummySessionRouter()

	response := mock.Request(t, http.StatusOK, router, "GET", "/active?collectors=foo", nil)

	mock.Validate(t, api.SessionsActive{}, response.Data)
}
