package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type jsonValidator struct {
	validator *validator.Validate
}

// New returns a new Validator for the echo webserver framework
func New() echo.Validator {
	v := &jsonValidator{
		validator: validator.New(),
	}

	return v
}

func (cv *jsonValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
