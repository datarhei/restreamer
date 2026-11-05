package api

import (
	"net/http"
	"strings"

	"github.com/datarhei/core/v16/http/api"
	"github.com/datarhei/core/v16/http/handler/util"
	"github.com/datarhei/core/v16/restream"
	"github.com/datarhei/core/v16/session"

	"github.com/labstack/echo/v4"
)

type WidgetConfig struct {
	Restream restream.Restreamer
	Registry session.RegistryReader
}

// The WidgetHandler type provides handlers for the widget API
type WidgetHandler struct {
	restream restream.Restreamer
	registry session.RegistryReader
}

// NewWidget return a new Widget type
func NewWidget(config WidgetConfig) *WidgetHandler {
	if config.Registry == nil {
		config.Registry, _ = session.New(session.Config{})
	}

	return &WidgetHandler{
		restream: config.Restream,
		registry: config.Registry,
	}
}

// Get returns minimal public statistics about a process
// @Summary Fetch minimal statistics about a process
// @Description Fetch minimal statistics about a process, which is not protected by any auth.
// @Tags v16.7.2
// @ID widget-3-get
// @Produce json
// @Param id path string true "ID of a process"
// @Success 200 {object} api.WidgetProcess
// @Failure 404 {object} api.Error
// @Router /api/v3/widget/process/{id} [get]
func (w *WidgetHandler) Get(c echo.Context) error {
	id := util.PathParam(c, "id")

	if w.restream == nil {
		return api.Err(http.StatusNotFound, "Unknown process ID")
	}

	process, err := w.restream.GetProcess(id)
	if err != nil {
		return api.Err(http.StatusNotFound, "Unknown process ID", "%s", err)
	}

	state, err := w.restream.GetProcessState(id)
	if err != nil {
		return api.Err(http.StatusNotFound, "Unknown process ID", "%s", err)
	}

	data := api.WidgetProcess{
		Uptime: int64(state.Duration),
	}

	if state.State != "running" {
		data.Uptime = 0
	}

	if w.registry == nil {
		return c.JSON(http.StatusOK, data)
	}

	collector := w.registry.Collector("hls")

	summary := collector.Summary()

	for _, session := range summary.Active {
		if !strings.HasPrefix(session.Reference, process.Reference) {
			continue
		}

		data.CurrentSessions++
	}

	for reference, s := range summary.Summary.References {
		if !strings.HasPrefix(reference, process.Reference) {
			continue
		}

		data.TotalSessions += s.TotalSessions
	}

	return c.JSON(http.StatusOK, data)
}
