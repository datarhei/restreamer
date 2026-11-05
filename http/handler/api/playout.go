package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/datarhei/core/v16/http/api"
	"github.com/datarhei/core/v16/http/handler/util"
	"github.com/datarhei/core/v16/playout"
	"github.com/datarhei/core/v16/restream"

	"github.com/labstack/echo/v4"
)

// The PlayoutHandler type provides handlers for accessing the playout API of a process
type PlayoutHandler struct {
	restream restream.Restreamer
}

// NewPlayout returns a new Playout type. You have to provide a Restreamer instance.
func NewPlayout(restream restream.Restreamer) *PlayoutHandler {
	return &PlayoutHandler{
		restream: restream,
	}
}

// Status return the current playout status
// @Summary Get the current playout status
// @Description Get the current playout status of an input of a process
// @Tags v16.7.2
// @ID process-3-playout-status
// @Produce json
// @Param id path string true "Process ID"
// @Param inputid path string true "Process Input ID"
// @Success 200 {object} api.PlayoutStatus
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/process/{id}/playout/{inputid}/status [get]
func (h *PlayoutHandler) Status(c echo.Context) error {
	id := util.PathParam(c, "id")
	inputid := util.PathParam(c, "inputid")

	addr, err := h.restream.GetPlayout(id, inputid)
	if err != nil {
		return api.Err(http.StatusNotFound, "Unknown process or input", "%s", err)
	}

	path := "/v1/status"

	response, err := h.request(http.MethodGet, addr, path, "", nil)
	if err != nil {
		return api.Err(http.StatusInternalServerError, "", "%s", err)
	}

	defer response.Body.Close()

	// Read the whole response
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return api.Err(http.StatusInternalServerError, "", "%s", err)
	}

	if response.StatusCode == http.StatusOK {
		status := playout.Status{}

		err := json.Unmarshal(data, &status)
		if err != nil {
			return api.Err(http.StatusInternalServerError, "", "%s", err)
		}

		apistatus := api.PlayoutStatus{}
		apistatus.Unmarshal(status)

		return c.JSON(http.StatusOK, apistatus)
	}

	return c.Blob(response.StatusCode, response.Header.Get("content-type"), data)
}

// Keyframe returns the last keyframe
// @Summary Get the last keyframe
// @Description Get the last keyframe of an input of a process. The extension of the name determines the return type.
// @Tags v16.7.2
// @ID process-3-playout-keyframe
// @Produce image/jpeg
// @Produce image/png
// @Produce json
// @Param id path string true "Process ID"
// @Param inputid path string true "Process Input ID"
// @Param name path string true "Any filename with an extension of .jpg or .png"
// @Success 200 {file} byte
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/process/{id}/playout/{inputid}/keyframe/{name} [get]
func (h *PlayoutHandler) Keyframe(c echo.Context) error {
	id := util.PathParam(c, "id")
	inputid := util.PathParam(c, "inputid")
	name := util.PathWildcardParam(c)

	addr, err := h.restream.GetPlayout(id, inputid)
	if err != nil {
		return api.Err(http.StatusNotFound, "Unknown process or input", "%s", err)
	}

	path := "/v1/keyframe/last."

	if strings.HasSuffix(name, ".png") {
		path = path + "png"
	} else {
		path = path + "jpg"
	}

	response, err := h.request(http.MethodGet, addr, path, "", nil)
	if err != nil {
		return api.Err(http.StatusInternalServerError, "", "%s", err)
	}

	defer response.Body.Close()

	// Read the whole response
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return api.Err(http.StatusInternalServerError, "", "%s", err)
	}

	return c.Blob(response.StatusCode, response.Header.Get("content-type"), data)
}

// EncodeErrorframe encodes the errorframe
// @Summary Encode the errorframe
// @Description Immediately encode the errorframe (if available and looping)
// @Tags v16.7.2
// @ID process-3-playout-errorframencode
// @Produce text/plain
// @Produce json
// @Param id path string true "Process ID"
// @Param inputid path string true "Process Input ID"
// @Success 204 {string} string
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/process/{id}/playout/{inputid}/errorframe/encode [get]
func (h *PlayoutHandler) EncodeErrorframe(c echo.Context) error {
	id := util.PathParam(c, "id")
	inputid := util.PathParam(c, "inputid")

	addr, err := h.restream.GetPlayout(id, inputid)
	if err != nil {
		return api.Err(http.StatusNotFound, "Unknown process or input", "%s", err)
	}

	path := "/v1/errorframe/encode"

	response, err := h.request(http.MethodGet, addr, path, "", nil)
	if err != nil {
		return api.Err(http.StatusInternalServerError, "", "%s", err)
	}

	defer response.Body.Close()

	// Read the whole response
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return api.Err(http.StatusInternalServerError, "", "%s", err)
	}

	return c.Blob(response.StatusCode, response.Header.Get("content-type"), data)
}

// SetErrorframe sets an errorframe
// @Summary Upload an error frame
// @Description Upload an error frame which will be encoded immediately
// @Tags v16.7.2
// @ID process-3-playout-errorframe
// @Produce text/plain
// @Produce json
// @Accept application/octet-stream
// @Param id path string true "Process ID"
// @Param inputid path string true "Process Input ID"
// @Param name path string true "Any filename with a suitable extension"
// @Param image body []byte true "Image to be used a error frame"
// @Success 204 {string} string
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/process/{id}/playout/{inputid}/errorframe/{name} [post]
func (h *PlayoutHandler) SetErrorframe(c echo.Context) error {
	id := util.PathParam(c, "id")
	inputid := util.PathParam(c, "inputid")

	addr, err := h.restream.GetPlayout(id, inputid)
	if err != nil {
		return api.Err(http.StatusNotFound, "Unknown process or input", "%s", err)
	}

	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return api.Err(http.StatusBadRequest, "Failed to read request body", "%s", err)
	}

	path := "/v1/errorframe.jpg"

	response, err := h.request(http.MethodPut, addr, path, "application/octet-stream", data)
	if err != nil {
		return api.Err(http.StatusInternalServerError, "", "%s", err)
	}

	defer response.Body.Close()

	// Read the whole response
	data, err = io.ReadAll(response.Body)
	if err != nil {
		return api.Err(http.StatusInternalServerError, "", "%s", err)
	}

	return c.Blob(response.StatusCode, response.Header.Get("content-type"), data)
}

// ReopenInput closes the current input stream
// @Summary Close the current input stream
// @Description Close the current input stream such that it will be automatically re-opened
// @Tags v16.7.2
// @ID process-3-playout-reopen-input
// @Produce plain
// @Param id path string true "Process ID"
// @Param inputid path string true "Process Input ID"
// @Success 200 {string} string
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/process/{id}/playout/{inputid}/reopen [get]
func (h *PlayoutHandler) ReopenInput(c echo.Context) error {
	id := util.PathParam(c, "id")
	inputid := util.PathParam(c, "inputid")

	addr, err := h.restream.GetPlayout(id, inputid)
	if err != nil {
		return api.Err(http.StatusNotFound, "Unknown process or input", "%s", err)
	}

	path := "/v1/reopen"

	response, err := h.request(http.MethodGet, addr, path, "", nil)
	if err != nil {
		return api.Err(http.StatusInternalServerError, "", "%s", err)
	}

	defer response.Body.Close()

	// Read the whole response
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return api.Err(http.StatusInternalServerError, "", "%s", err)
	}

	return c.Blob(response.StatusCode, response.Header.Get("content-type"), data)
}

// SetStream replaces the current stream
// @Summary Switch to a new stream
// @Description Replace the current stream with the one from the given URL. The switch will only happen if the stream parameters match.
// @Tags v16.7.2
// @ID process-3-playout-stream
// @Produce text/plain
// @Produce json
// @Accept text/plain
// @Param id path string true "Process ID"
// @Param inputid path string true "Process Input ID"
// @Param url body string true "URL of the new stream"
// @Success 204 {string} string
// @Failure 404 {object} api.Error
// @Failure 500 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/process/{id}/playout/{inputid}/stream [put]
func (h *PlayoutHandler) SetStream(c echo.Context) error {
	id := util.PathParam(c, "id")
	inputid := util.PathParam(c, "inputid")

	addr, err := h.restream.GetPlayout(id, inputid)
	if err != nil {
		return api.Err(http.StatusNotFound, "Unknown process or input", "%s", err)
	}

	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return api.Err(http.StatusBadRequest, "Failed to read request body", "%s", err)
	}

	path := "/v1/stream"

	response, err := h.request(http.MethodPut, addr, path, "text/plain", data)
	if err != nil {
		return api.Err(http.StatusInternalServerError, "", "%s", err)
	}

	defer response.Body.Close()

	// Read the whole response
	data, err = io.ReadAll(response.Body)
	if err != nil {
		return api.Err(http.StatusInternalServerError, "", "%s", err)
	}

	return c.Blob(response.StatusCode, response.Header.Get("content-type"), data)
}

func (h *PlayoutHandler) request(method, addr, path, contentType string, data []byte) (*http.Response, error) {
	endpoint := "http://" + addr + path

	body := bytes.NewBuffer(data)

	request, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", contentType)

	// Submit the request
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
