package errorhandler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/datarhei/core/v16/http/api"

	"github.com/labstack/echo/v4"
)

// HTTPErrorHandler is a genral handler for echo handler errors
func HTTPErrorHandler(err error, c echo.Context) {
	var code int = 0
	var details []string
	message := ""

	if he, ok := err.(api.Error); ok {
		code = he.Code
		message = he.Message
		details = he.Details
	} else if he, ok := err.(*echo.HTTPError); ok {
		if he.Internal != nil {
			if herr, ok := he.Internal.(*echo.HTTPError); ok {
				he = herr
			}
		}

		code = he.Code
		message = http.StatusText(he.Code)
		if len(message) == 0 {
			switch code {
			case 509:
				message = "Bandwith limit exceeded"
			default:
			}
		}
		details = strings.Split(fmt.Sprintf("%v", he.Message), "\n")
	} else {
		code = http.StatusInternalServerError
		message = http.StatusText(http.StatusInternalServerError)
		details = strings.Split(fmt.Sprintf("%s", err), "\n")
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			c.NoContent(code)
		} else {
			c.JSON(code, api.Error{
				Code:    code,
				Message: message,
				Details: details,
			})
		}
	}
}
