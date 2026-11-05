package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// The PrometheusHandler type provides a handler function for reading the prometheus metrics
type PrometheusHandler struct {
	handler http.Handler
}

// NewPrometheus returns a new Prometheus type. You have to provide a HTTP.Handler
func NewPrometheus(handler http.Handler) *PrometheusHandler {
	return &PrometheusHandler{
		handler: handler,
	}
}

// Metrics godoc
// @Summary Prometheus metrics
// @Description Prometheus metrics
// @ID metrics
// @Produce text/plain
// @Success 200 {string} string
// @Router /metrics [get]
func (m *PrometheusHandler) Metrics(c echo.Context) error {
	m.handler.ServeHTTP(c.Response(), c.Request())

	return nil
}
