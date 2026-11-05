package api

import (
	"net/http"
	"testing"

	"github.com/datarhei/core/v16/http/api"
	"github.com/datarhei/core/v16/http/mock"
	"github.com/stretchr/testify/require"

	"github.com/labstack/echo/v4"
)

func getDummyAboutRouter() (*echo.Echo, error) {
	router := mock.DummyEcho()

	rs, err := mock.DummyRestreamer("../../mock")
	if err != nil {
		return nil, err
	}

	handler := NewAbout(rs, []string{})

	router.Add("GET", "/", handler.About)

	return router, nil
}

func TestAbout(t *testing.T) {
	router, err := getDummyAboutRouter()
	require.NoError(t, err)

	response := mock.Request(t, http.StatusOK, router, "GET", "/", nil)

	mock.Validate(t, &api.About{}, response.Data)
}
