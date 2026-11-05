// Package api implements various handlers for the API routes
package api

import (
	"github.com/labstack/echo/v4"
)

// Swagger returns swagger UI for this API
// @Summary Swagger UI for this API
// @Description Swagger UI for this API
// @ID swagger
// @Produce text/html
// @Success 200 {string} string ""
// @Router /api/swagger [get]
func Swagger(c echo.Context) error {
	return nil
}
