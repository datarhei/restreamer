package util

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/datarhei/core/v16/encoding/json"

	"github.com/labstack/echo/v4"
)

func ShouldBindJSONValidation(c echo.Context, obj interface{}, validate bool) error {
	req := c.Request()

	if req.ContentLength == 0 {
		return fmt.Errorf("request doesn't contain any content")
	}

	ctype := req.Header.Get(echo.HeaderContentType)

	if !strings.HasPrefix(ctype, echo.MIMEApplicationJSON) {
		return fmt.Errorf("request doesn't contain JSON content")
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, obj); err != nil {
		return json.FormatError(body, err)
	}

	if validate {
		return c.Validate(obj)
	}

	return nil
}

// ShouldBindJSON binds the body data of the request to the given object. An error is
// returned if the body data is not valid JSON or the validation of the unmarshalled
// data failed.
func ShouldBindJSON(c echo.Context, obj interface{}) error {
	return ShouldBindJSONValidation(c, obj, true)
}

func PathWildcardParam(c echo.Context) string {
	return "/" + PathParam(c, "*")
}

func PathParam(c echo.Context, name string) string {
	param := c.Param(name)

	param, err := url.PathUnescape(param)
	if err != nil {
		return ""
	}

	return param
}

func DefaultQuery(c echo.Context, name, defValue string) string {
	param := c.QueryParam(name)

	if len(param) == 0 {
		return defValue
	}

	param, err := url.QueryUnescape(param)
	if err != nil {
		return defValue
	}

	return param
}
