package api

import (
	"net/http"
	"strings"

	"github.com/datarhei/core/v16/http/handler/util"
	"github.com/datarhei/core/v16/log"

	"github.com/labstack/echo/v4"
)

// The LogHandler type provides handler functions for reading the application log
type LogHandler struct {
	buffer log.BufferWriter
}

// NewLog return a new Log type. You have to provide log buffer.
func NewLog(buffer log.BufferWriter) *LogHandler {
	l := &LogHandler{
		buffer: buffer,
	}

	if l.buffer == nil {
		l.buffer = log.NewBufferWriter(log.Lsilent, 1)
	}

	return l
}

// Log returns the last log lines of the Restreamer application
// @Summary Application log
// @Description Get the last log lines of the Restreamer application
// @Tags v16.7.2
// @ID log-3
// @Param format query string false "Format of the list of log events (*console, raw)"
// @Produce json
// @Success 200 {array} api.LogEvent "application log"
// @Success 200 {array} string "application log"
// @Security ApiKeyAuth
// @Router /api/v3/log [get]
func (p *LogHandler) Log(c echo.Context) error {
	format := util.DefaultQuery(c, "format", "console")

	events := p.buffer.Events()

	if format == "raw" {
		log := make([]map[string]interface{}, len(events))

		for i, e := range events {
			e.Data["ts"] = e.Time
			e.Data["component"] = e.Component

			if len(e.Caller) != 0 {
				e.Data["caller"] = e.Caller
			}

			if len(e.Message) != 0 {
				e.Data["message"] = e.Message
			}

			log[i] = e.Data
		}

		return c.JSON(http.StatusOK, log)
	}

	formatter := log.NewConsoleFormatter(false)

	log := make([]string, len(events))

	for i, e := range events {
		log[i] = strings.TrimSpace(formatter.String(e))
	}

	return c.JSON(http.StatusOK, log)
}
