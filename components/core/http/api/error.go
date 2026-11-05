package api

import (
	"fmt"
	"net/http"
	"strings"
)

// Error represents an error response of the API
type Error struct {
	Code    int      `json:"code" jsonschema:"required" format:"int"`
	Message string   `json:"message" jsonschema:""`
	Details []string `json:"details" jsonschema:""`
}

// Error returns the string representation of the error
func (e Error) Error() string {
	return fmt.Sprintf("code=%d, message=%s, details=%s", e.Code, e.Message, strings.Join(e.Details, " "))
}

// Err creates a new API error with the given HTTP status code. If message is empty, the default message
// for the given code is used. If the first entry in args is a string, it is interpreted as a format string
// for the remaining entries in args, that is used for fmt.Sprintf. Otherwise the args are ignored.
func Err(code int, message string, args ...interface{}) Error {
	if len(message) == 0 {
		message = http.StatusText(code)
	}

	e := Error{
		Code:    code,
		Message: message,
		Details: []string{},
	}

	if len(args) >= 1 {
		if format, ok := args[0].(string); ok {
			e.Details = strings.Split(fmt.Sprintf(format, args[1:]...), "\n")
		}
	}

	return e
}
