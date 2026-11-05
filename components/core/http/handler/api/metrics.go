package api

import (
	"net/http"
	"sort"
	"time"

	"github.com/datarhei/core/v16/http/api"
	"github.com/datarhei/core/v16/http/handler/util"
	"github.com/datarhei/core/v16/monitor"
	"github.com/datarhei/core/v16/monitor/metric"

	"github.com/labstack/echo/v4"
)

type MetricsConfig struct {
	Metrics monitor.HistoryReader
}

// The MetricsHandler type provides handlers for the widget API
type MetricsHandler struct {
	metrics monitor.HistoryReader
}

// NewWidget return a new Widget type
func NewMetrics(config MetricsConfig) *MetricsHandler {
	return &MetricsHandler{
		metrics: config.Metrics,
	}
}

// Describe the known metrics
// @Summary List all known metrics with their description and labels
// @Description List all known metrics with their description and labels
// @Tags v16.10.0
// @ID metrics-3-describe
// @Produce json
// @Success 200 {array} api.MetricsDescription
// @Security ApiKeyAuth
// @Router /api/v3/metrics [get]
func (r *MetricsHandler) Describe(c echo.Context) error {
	response := []api.MetricsDescription{}

	descriptors := r.metrics.Describe()

	for _, d := range descriptors {
		response = append(response, api.MetricsDescription{
			Name:        d.Name(),
			Description: d.Description(),
			Labels:      d.Labels(),
		})
	}

	sort.Slice(response, func(i, j int) bool {
		return response[i].Name < response[j].Name
	})

	return c.JSON(http.StatusOK, response)
}

// Query the collected metrics
// @Summary Query the collected metrics
// @Description Query the collected metrics
// @Tags v16.7.2
// @ID metrics-3-metrics
// @Accept json
// @Produce json
// @Param config body api.MetricsQuery true "Metrics query"
// @Success 200 {object} api.MetricsResponse
// @Failure 400 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/metrics [post]
func (r *MetricsHandler) Metrics(c echo.Context) error {
	var query api.MetricsQuery

	if err := util.ShouldBindJSON(c, &query); err != nil {
		return api.Err(http.StatusBadRequest, "Invalid JSON", "%s", err)
	}

	patterns := []metric.Pattern{}

	for _, m := range query.Metrics {
		labels := []string{}
		for k, v := range m.Labels {
			labels = append(labels, k, v)
		}

		pattern := metric.NewPattern(m.Name, labels...)

		patterns = append(patterns, pattern)
	}

	response := api.MetricsResponse{}

	var data []monitor.HistoryMetrics

	if query.Timerange == 0 {
		// current data
		data = []monitor.HistoryMetrics{
			{
				TS:      time.Now(),
				Metrics: r.metrics.Collect(patterns),
			},
		}
	} else {
		// historic data
		data = r.metrics.History(time.Second*time.Duration(query.Timerange), time.Second*time.Duration(query.Interval), patterns)
	}

	timerange, interval := r.metrics.Resolution()

	response.Unmarshal(data, timerange, interval)

	return c.JSON(http.StatusOK, response)
}
