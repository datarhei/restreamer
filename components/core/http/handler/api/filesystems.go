package api

import (
	"net/http"

	"github.com/datarhei/core/v16/http/api"
	"github.com/datarhei/core/v16/http/handler"
	"github.com/datarhei/core/v16/http/handler/util"

	"github.com/labstack/echo/v4"
)

type FSConfig struct {
	Type       string
	Mountpoint string
	Handler    *handler.FSHandler
}

// The FSHandler type provides handlers for manipulating a filesystem
type FSHandler struct {
	filesystems map[string]FSConfig
}

// NewFS return a new FSHanlder type. You have to provide a filesystem to act on.
func NewFS(filesystems map[string]FSConfig) *FSHandler {
	return &FSHandler{
		filesystems: filesystems,
	}
}

// GetFileAPI returns the file at the given path
// @Summary Fetch a file from a filesystem
// @Description Fetch a file from a filesystem
// @Tags v16.7.2
// @ID filesystem-3-get-file
// @Produce application/data
// @Produce json
// @Param storage path string true "Name of the filesystem"
// @Param filepath path string true "Path to file"
// @Success 200 {file} byte
// @Success 301 {string} string
// @Failure 404 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/fs/{storage}/{filepath} [get]
func (h *FSHandler) GetFile(c echo.Context) error {
	name := util.PathParam(c, "name")

	config, ok := h.filesystems[name]
	if !ok {
		return api.Err(http.StatusNotFound, "File not found", "unknown filesystem: %s", name)
	}

	return config.Handler.GetFile(c)
}

// PutFileAPI adds or overwrites a file at the given path
// @Summary Add a file to a filesystem
// @Description Writes or overwrites a file on a filesystem
// @Tags v16.7.2
// @ID filesystem-3-put-file
// @Accept application/data
// @Produce text/plain
// @Produce json
// @Param storage path string true "Name of the filesystem"
// @Param filepath path string true "Path to file"
// @Param data body []byte true "File data"
// @Success 201 {string} string
// @Success 204 {string} string
// @Failure 507 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/fs/{storage}/{filepath} [put]
func (h *FSHandler) PutFile(c echo.Context) error {
	name := util.PathParam(c, "name")

	config, ok := h.filesystems[name]
	if !ok {
		return api.Err(http.StatusNotFound, "File not found", "unknown filesystem: %s", name)
	}

	return config.Handler.PutFile(c)
}

// DeleteFileAPI removes a file from a filesystem
// @Summary Remove a file from a filesystem
// @Description Remove a file from a filesystem
// @Tags v16.7.2
// @ID filesystem-3-delete-file
// @Produce text/plain
// @Param storage path string true "Name of the filesystem"
// @Param filepath path string true "Path to file"
// @Success 200 {string} string
// @Failure 404 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/fs/{storage}/{filepath} [delete]
func (h *FSHandler) DeleteFile(c echo.Context) error {
	name := util.PathParam(c, "name")

	config, ok := h.filesystems[name]
	if !ok {
		return api.Err(http.StatusNotFound, "File not found", "unknown filesystem: %s", name)
	}

	return config.Handler.DeleteFile(c)
}

// ListFiles lists all files on a filesystem
// @Summary List all files on a filesystem
// @Description List all files on a filesystem. The listing can be ordered by name, size, or date of last modification in ascending or descending order.
// @Tags v16.7.2
// @ID filesystem-3-list-files
// @Produce json
// @Param storage path string true "Name of the filesystem"
// @Param glob query string false "glob pattern for file names"
// @Param sort query string false "none, name, size, lastmod"
// @Param order query string false "asc, desc"
// @Success 200 {array} api.FileInfo
// @Security ApiKeyAuth
// @Router /api/v3/fs/{storage} [get]
func (h *FSHandler) ListFiles(c echo.Context) error {
	name := util.PathParam(c, "name")

	config, ok := h.filesystems[name]
	if !ok {
		return api.Err(http.StatusNotFound, "File not found", "unknown filesystem: %s", name)
	}

	return config.Handler.ListFiles(c)
}

// List lists all registered filesystems
// @Summary List all registered filesystems
// @Description Listall registered filesystems
// @Tags v16.12.0
// @ID filesystem-3-list
// @Produce json
// @Success 200 {array} api.FilesystemInfo
// @Security ApiKeyAuth
// @Router /api/v3/fs [get]
func (h *FSHandler) List(c echo.Context) error {
	fss := []api.FilesystemInfo{}

	for name, config := range h.filesystems {
		fss = append(fss, api.FilesystemInfo{
			Name:  name,
			Type:  config.Type,
			Mount: config.Mountpoint,
		})
	}

	return c.JSON(http.StatusOK, fss)
}
