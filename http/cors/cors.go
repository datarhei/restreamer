// Package cors implements a validator for CORS origins
package cors

import (
	"errors"
	"strings"
)

// Copied from github.com/gin-contrib/cors

// DefaultSchemas is a list of default allowed schemas for CORS origins
var DefaultSchemas = []string{
	"http://",
	"https://",
}

// Validate checks a list of origins if the comply with the allowed origins
func Validate(origins []string) error {
	for _, origin := range origins {
		if !strings.Contains(origin, "*") && !validateAllowedSchemas(origin) {
			return errors.New("bad origin: origins must contain '*' or include " + strings.Join(getAllowedSchemas(), ", or "))
		}
	}
	return nil
}

func validateAllowedSchemas(origin string) bool {
	allowedSchemas := getAllowedSchemas()

	for _, schema := range allowedSchemas {
		if strings.HasPrefix(origin, schema) {
			return true
		}
	}

	return false
}

func getAllowedSchemas() []string {
	allowedSchemas := DefaultSchemas

	return allowedSchemas
}
