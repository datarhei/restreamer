package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Null is a handler that always returns NoContent
func Null(c echo.Context) error {
	return c.String(http.StatusNoContent, "")
}
