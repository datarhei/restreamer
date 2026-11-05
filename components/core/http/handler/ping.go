package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// The PingHandler type provides a handler for a ping requesr
type PingHandler struct{}

// NewPing returns a neww Ping type.
func NewPing() *PingHandler {
	return &PingHandler{}
}

// Ping returns pong
// @Summary Liveliness check
// @Description Liveliness check
// @ID ping
// @Produce text/plain
// @Success 200 {string} string "pong"
// @Router /ping [get]
func (p *PingHandler) Ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
