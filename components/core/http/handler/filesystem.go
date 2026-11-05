package handler

import (
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"github.com/datarhei/core/v16/http/api"
	"github.com/datarhei/core/v16/http/fs"
	"github.com/datarhei/core/v16/http/handler/util"

	"github.com/labstack/echo/v4"
)

// The FSHandler type provides handlers for manipulating a filesystem
type FSHandler struct {
	fs fs.FS
}

// NewFS return a new FSHandler type. You have to provide a filesystem to act on.
func NewFS(fs fs.FS) *FSHandler {
	return &FSHandler{
		fs: fs,
	}
}

func (h *FSHandler) GetFile(c echo.Context) error {
	path := util.PathWildcardParam(c)

	mimeType := c.Response().Header().Get(echo.HeaderContentType)
	c.Response().Header().Del(echo.HeaderContentType)

	file := h.fs.Filesystem.Open(path)
	if file == nil {
		return api.Err(http.StatusNotFound, "File not found", path)
	}

	stat, _ := file.Stat()

	if len(h.fs.DefaultFile) != 0 {
		if stat.IsDir() {
			path = filepath.Join(path, h.fs.DefaultFile)

			file.Close()

			file = h.fs.Filesystem.Open(path)
			if file == nil {
				return api.Err(http.StatusNotFound, "File not found", path)
			}

			stat, _ = file.Stat()
		}
	}

	defer file.Close()

	c.Response().Header().Set("Last-Modified", stat.ModTime().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))

	if path, ok := stat.IsLink(); ok {
		path = filepath.Clean("/" + path)

		if path[0] == '/' {
			path = path[1:]
		}

		return c.Redirect(http.StatusMovedPermanently, path)
	}

	c.Response().Header().Set(echo.HeaderContentType, mimeType)

	if c.Request().Method == "HEAD" {
		return c.Blob(http.StatusOK, "application/data", nil)
	}

	return c.Stream(http.StatusOK, "application/data", file)
}

func (h *FSHandler) PutFile(c echo.Context) error {
	path := util.PathWildcardParam(c)

	c.Response().Header().Del(echo.HeaderContentType)

	req := c.Request()

	_, created, err := h.fs.Filesystem.WriteFileReader(path, req.Body)
	if err != nil {
		return api.Err(http.StatusBadRequest, "Bad request", "%s", err)
	}

	if h.fs.Cache != nil {
		h.fs.Cache.Delete(path)

		if len(h.fs.DefaultFile) != 0 {
			if strings.HasSuffix(path, "/"+h.fs.DefaultFile) {
				path := strings.TrimSuffix(path, h.fs.DefaultFile)
				h.fs.Cache.Delete(path)
			}
		}
	}

	c.Response().Header().Set("Content-Location", req.URL.RequestURI())

	if created {
		return c.String(http.StatusCreated, "")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *FSHandler) DeleteFile(c echo.Context) error {
	path := util.PathWildcardParam(c)

	c.Response().Header().Del(echo.HeaderContentType)

	size := h.fs.Filesystem.Remove(path)

	if h.fs.Cache != nil {
		h.fs.Cache.Delete(path)

		if len(h.fs.DefaultFile) != 0 {
			if strings.HasSuffix(path, "/"+h.fs.DefaultFile) {
				path := strings.TrimSuffix(path, h.fs.DefaultFile)
				h.fs.Cache.Delete(path)
			}
		}
	}

	if size < 0 {
		return api.Err(http.StatusNotFound, "File not found", path)
	}

	return c.String(http.StatusOK, "Deleted: "+path)
}

func (h *FSHandler) ListFiles(c echo.Context) error {
	pattern := util.DefaultQuery(c, "glob", "")
	sortby := util.DefaultQuery(c, "sort", "none")
	order := util.DefaultQuery(c, "order", "asc")

	files := h.fs.Filesystem.List("/", pattern)

	var sortFunc func(i, j int) bool

	switch sortby {
	case "name":
		if order == "desc" {
			sortFunc = func(i, j int) bool { return files[i].Name() > files[j].Name() }
		} else {
			sortFunc = func(i, j int) bool { return files[i].Name() < files[j].Name() }
		}
	case "size":
		if order == "desc" {
			sortFunc = func(i, j int) bool { return files[i].Size() > files[j].Size() }
		} else {
			sortFunc = func(i, j int) bool { return files[i].Size() < files[j].Size() }
		}
	default:
		if order == "asc" {
			sortFunc = func(i, j int) bool { return files[i].ModTime().Before(files[j].ModTime()) }
		} else {
			sortFunc = func(i, j int) bool { return files[i].ModTime().After(files[j].ModTime()) }
		}
	}

	sort.Slice(files, sortFunc)

	var fileinfos []api.FileInfo = make([]api.FileInfo, len(files))

	for i, f := range files {
		fileinfos[i] = api.FileInfo{
			Name:    f.Name(),
			Size:    f.Size(),
			LastMod: f.ModTime().Unix(),
		}
	}

	return c.JSON(http.StatusOK, fileinfos)
}
