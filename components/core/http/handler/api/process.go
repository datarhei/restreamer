package api

import (
	"net/http"
	"strings"

	"github.com/datarhei/core/v16/http/api"
	"github.com/datarhei/core/v16/http/handler/util"
	"github.com/datarhei/core/v16/restream"

	"github.com/labstack/echo/v4"
	"github.com/lithammer/shortuuid/v4"
)

// The RestreamHandler type provides functions to interact with a Restreamer instance
type RestreamHandler struct {
	restream restream.Restreamer
}

// NewRestream return a new Restream type. You have to provide a valid Restreamer instance.
func NewRestream(restream restream.Restreamer) *RestreamHandler {
	return &RestreamHandler{
		restream: restream,
	}
}

// Add adds a new process
// @Summary Add a new process
// @Description Add a new FFmpeg process
// @Tags v16.7.2
// @ID process-3-add
// @Accept json
// @Produce json
// @Param config body api.ProcessConfig true "Process config"
// @Success 200 {object} api.ProcessConfig
// @Failure 400 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/process [post]
func (h *RestreamHandler) Add(c echo.Context) error {
	process := api.ProcessConfig{
		ID:        shortuuid.New(),
		Type:      "ffmpeg",
		Autostart: true,
	}

	if err := util.ShouldBindJSON(c, &process); err != nil {
		return api.Err(http.StatusBadRequest, "Invalid JSON", "%s", err)
	}

	if process.Type != "ffmpeg" {
		return api.Err(http.StatusBadRequest, "Unsupported process type", "Supported process types are: ffmpeg")
	}

	if len(process.Input) == 0 || len(process.Output) == 0 {
		return api.Err(http.StatusBadRequest, "At least one input and one output need to be defined")
	}

	config := process.Marshal()

	if err := h.restream.AddProcess(config); err != nil {
		return api.Err(http.StatusBadRequest, "Invalid process config", "%s", err.Error())
	}

	p, _ := h.getProcess(config.ID, "config")

	return c.JSON(http.StatusOK, p.Config)
}

// GetAll returns all known processes
// @Summary List all known processes
// @Description List all known processes. Use the query parameter to filter the listed processes.
// @Tags v16.7.2
// @ID process-3-get-all
// @Produce json
// @Param filter query string false "Comma separated list of fields (config, state, report, metadata) that will be part of the output. If empty, all fields will be part of the output."
// @Param reference query string false "Return only these process that have this reference value. If empty, the reference will be ignored."
// @Param id query string false "Comma separated list of process ids to list. Overrides the reference. If empty all IDs will be returned."
// @Param idpattern query string false "Glob pattern for process IDs. If empty all IDs will be returned. Intersected with results from refpattern."
// @Param refpattern query string false "Glob pattern for process references. If empty all IDs will be returned. Intersected with results from idpattern."
// @Success 200 {array} api.Process
// @Security ApiKeyAuth
// @Router /api/v3/process [get]
func (h *RestreamHandler) GetAll(c echo.Context) error {
	filter := util.DefaultQuery(c, "filter", "")
	reference := util.DefaultQuery(c, "reference", "")
	wantids := strings.FieldsFunc(util.DefaultQuery(c, "id", ""), func(r rune) bool {
		return r == rune(',')
	})
	idpattern := util.DefaultQuery(c, "idpattern", "")
	refpattern := util.DefaultQuery(c, "refpattern", "")

	ids := h.restream.GetProcessIDs(idpattern, refpattern)

	processes := []api.Process{}

	if len(wantids) == 0 || len(reference) != 0 {
		for _, id := range ids {
			if p, err := h.getProcess(id, filter); err == nil {
				if len(reference) != 0 && p.Reference != reference {
					continue
				}
				processes = append(processes, p)
			}
		}
	} else {
		for _, id := range ids {
			for _, wantid := range wantids {
				if wantid == id {
					if p, err := h.getProcess(id, filter); err == nil {
						processes = append(processes, p)
					}
				}
			}
		}
	}

	return c.JSON(http.StatusOK, processes)
}

// Get returns the process with the given ID
// @Summary List a process by its ID
// @Description List a process by its ID. Use the filter parameter to specifiy the level of detail of the output.
// @Tags v16.7.2
// @ID process-3-get
// @Produce json
// @Param id path string true "Process ID"
// @Param filter query string false "Comma separated list of fields (config, state, report, metadata) to be part of the output. If empty, all fields will be part of the output"
// @Success 200 {object} api.Process
// @Failure 404 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/process/{id} [get]
func (h *RestreamHandler) Get(c echo.Context) error {
	id := util.PathParam(c, "id")
	filter := util.DefaultQuery(c, "filter", "")

	p, err := h.getProcess(id, filter)
	if err != nil {
		return api.Err(http.StatusNotFound, "Unknown process ID", "%s", err)
	}

	return c.JSON(http.StatusOK, p)
}

// Delete deletes the process with the given ID
// @Summary Delete a process by its ID
// @Description Delete a process by its ID
// @Tags v16.7.2
// @ID process-3-delete
// @Produce json
// @Param id path string true "Process ID"
// @Success 200 {string} string
// @Failure 404 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/process/{id} [delete]
func (h *RestreamHandler) Delete(c echo.Context) error {
	id := util.PathParam(c, "id")

	if err := h.restream.StopProcess(id); err != nil {
		return api.Err(http.StatusNotFound, "Unknown process ID", "%s", err)
	}

	if err := h.restream.DeleteProcess(id); err != nil {
		return api.Err(http.StatusInternalServerError, "Process can't be deleted", "%s", err)
	}

	return c.JSON(http.StatusOK, "OK")
}

// Update replaces an existing process
// @Summary Replace an existing process
// @Description Replace an existing process.
// @Tags v16.7.2
// @ID process-3-update
// @Accept json
// @Produce json
// @Param id path string true "Process ID"
// @Param config body api.ProcessConfig true "Process config"
// @Success 200 {object} api.ProcessConfig
// @Failure 400 {object} api.Error
// @Failure 404 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/process/{id} [put]
func (h *RestreamHandler) Update(c echo.Context) error {
	id := util.PathParam(c, "id")

	process := api.ProcessConfig{
		ID:        id,
		Type:      "ffmpeg",
		Autostart: true,
	}

	current, err := h.restream.GetProcess(id)
	if err != nil {
		return api.Err(http.StatusNotFound, "Process not found", "%s", id)
	}

	// Prefill the config with the current values
	process.Unmarshal(current.Config)

	if err := util.ShouldBindJSON(c, &process); err != nil {
		return api.Err(http.StatusBadRequest, "Invalid JSON", "%s", err)
	}

	config := process.Marshal()

	if err := h.restream.UpdateProcess(id, config); err != nil {
		if err == restream.ErrUnknownProcess {
			return api.Err(http.StatusNotFound, "Process not found", "%s", id)
		}

		return api.Err(http.StatusBadRequest, "Process can't be updated", "%s", err)
	}

	p, _ := h.getProcess(config.ID, "config")

	return c.JSON(http.StatusOK, p.Config)
}

// Command issues a command to a process
// @Summary Issue a command to a process
// @Description Issue a command to a process: start, stop, reload, restart
// @Tags v16.7.2
// @ID process-3-command
// @Accept json
// @Produce json
// @Param id path string true "Process ID"
// @Param command body api.Command true "Process command"
// @Success 200 {string} string
// @Failure 400 {object} api.Error
// @Failure 404 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/process/{id}/command [put]
func (h *RestreamHandler) Command(c echo.Context) error {
	id := util.PathParam(c, "id")

	var command api.Command

	if err := util.ShouldBindJSON(c, &command); err != nil {
		return api.Err(http.StatusBadRequest, "Invalid JSON", "%s", err)
	}

	var err error
	if command.Command == "start" {
		err = h.restream.StartProcess(id)
	} else if command.Command == "stop" {
		err = h.restream.StopProcess(id)
	} else if command.Command == "restart" {
		err = h.restream.RestartProcess(id)
	} else if command.Command == "reload" {
		err = h.restream.ReloadProcess(id)
	} else {
		return api.Err(http.StatusBadRequest, "Unknown command provided", "Known commands are: start, stop, reload, restart")
	}

	if err != nil {
		return api.Err(http.StatusBadRequest, "Command failed", "%s", err)
	}

	return c.JSON(http.StatusOK, "OK")
}

// GetConfig returns the configuration of a process
// @Summary Get the configuration of a process
// @Description Get the configuration of a process. This is the configuration as provided by Add or Update.
// @Tags v16.7.2
// @ID process-3-get-config
// @Produce json
// @Param id path string true "Process ID"
// @Success 200 {object} api.ProcessConfig
// @Failure 404 {object} api.Error
// @Failure 400 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/process/{id}/config [get]
func (h *RestreamHandler) GetConfig(c echo.Context) error {
	id := util.PathParam(c, "id")

	p, err := h.restream.GetProcess(id)
	if err != nil {
		return api.Err(http.StatusNotFound, "Unknown process ID", "%s", err)
	}

	config := api.ProcessConfig{}
	config.Unmarshal(p.Config)

	return c.JSON(http.StatusOK, config)
}

// GetState returns the current state of a process
// @Summary Get the state of a process
// @Description Get the state and progress data of a process.
// @Tags v16.7.2
// @ID process-3-get-state
// @Produce json
// @Param id path string true "Process ID"
// @Success 200 {object} api.ProcessState
// @Failure 404 {object} api.Error
// @Failure 400 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/process/{id}/state [get]
func (h *RestreamHandler) GetState(c echo.Context) error {
	id := util.PathParam(c, "id")

	s, err := h.restream.GetProcessState(id)
	if err != nil {
		return api.Err(http.StatusNotFound, "Unknown process ID", "%s", err)
	}

	state := api.ProcessState{}
	state.Unmarshal(s)

	return c.JSON(http.StatusOK, state)
}

// GetReport return the current log and the log history of a process
// @Summary Get the logs of a process
// @Description Get the logs and the log history of a process.
// @Tags v16.7.2
// @ID process-3-get-report
// @Produce json
// @Param id path string true "Process ID"
// @Success 200 {object} api.ProcessReport
// @Failure 404 {object} api.Error
// @Failure 400 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/process/{id}/report [get]
func (h *RestreamHandler) GetReport(c echo.Context) error {
	id := util.PathParam(c, "id")

	l, err := h.restream.GetProcessLog(id)
	if err != nil {
		return api.Err(http.StatusNotFound, "Unknown process ID", "%s", err)
	}

	report := api.ProcessReport{}
	report.Unmarshal(l)

	return c.JSON(http.StatusOK, report)
}

// Probe probes a process
// @Summary Probe a process
// @Description Probe an existing process to get a detailed stream information on the inputs.
// @Tags v16.7.2
// @ID process-3-probe
// @Produce json
// @Param id path string true "Process ID"
// @Success 200 {object} api.Probe
// @Security ApiKeyAuth
// @Router /api/v3/process/{id}/probe [get]
func (h *RestreamHandler) Probe(c echo.Context) error {
	id := util.PathParam(c, "id")

	probe := h.restream.Probe(id)

	apiprobe := api.Probe{}
	apiprobe.Unmarshal(&probe)

	return c.JSON(http.StatusOK, apiprobe)
}

// Skills returns the detected FFmpeg capabilities
// @Summary FFmpeg capabilities
// @Description List all detected FFmpeg capabilities.
// @Tags v16.7.2
// @ID skills-3
// @Produce json
// @Success 200 {object} api.Skills
// @Security ApiKeyAuth
// @Router /api/v3/skills [get]
func (h *RestreamHandler) Skills(c echo.Context) error {
	skills := h.restream.Skills()

	apiskills := api.Skills{}
	apiskills.Unmarshal(skills)

	return c.JSON(http.StatusOK, apiskills)
}

// ReloadSkills will refresh the FFmpeg capabilities
// @Summary Refresh FFmpeg capabilities
// @Description Refresh the available FFmpeg capabilities.
// @Tags v16.7.2
// @ID skills-3-reload
// @Produce json
// @Success 200 {object} api.Skills
// @Security ApiKeyAuth
// @Router /api/v3/skills/reload [get]
func (h *RestreamHandler) ReloadSkills(c echo.Context) error {
	h.restream.ReloadSkills()
	skills := h.restream.Skills()

	apiskills := api.Skills{}
	apiskills.Unmarshal(skills)

	return c.JSON(http.StatusOK, apiskills)
}

// GetProcessMetadata returns the metadata stored with a process
// @Summary Retrieve JSON metadata stored with a process under a key
// @Description Retrieve the previously stored JSON metadata under the given key. If the key is empty, all metadata will be returned.
// @Tags v16.7.2
// @ID process-3-get-process-metadata
// @Produce json
// @Param id path string true "Process ID"
// @Param key path string true "Key for data store"
// @Success 200 {object} api.Metadata
// @Failure 404 {object} api.Error
// @Failure 400 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/process/{id}/metadata/{key} [get]
func (h *RestreamHandler) GetProcessMetadata(c echo.Context) error {
	id := util.PathParam(c, "id")
	key := util.PathParam(c, "key")

	data, err := h.restream.GetProcessMetadata(id, key)
	if err != nil {
		return api.Err(http.StatusNotFound, "Unknown process ID", "%s", err)
	}

	return c.JSON(http.StatusOK, data)
}

// SetProcessMetadata stores metadata with a process
// @Summary Add JSON metadata with a process under the given key
// @Description Add arbitrary JSON metadata under the given key. If the key exists, all already stored metadata with this key will be overwritten. If the key doesn't exist, it will be created.
// @Tags v16.7.2
// @ID process-3-set-process-metadata
// @Produce json
// @Param id path string true "Process ID"
// @Param key path string true "Key for data store"
// @Param data body api.Metadata true "Arbitrary JSON data. The null value will remove the key and its contents"
// @Success 200 {object} api.Metadata
// @Failure 404 {object} api.Error
// @Failure 400 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/process/{id}/metadata/{key} [put]
func (h *RestreamHandler) SetProcessMetadata(c echo.Context) error {
	id := util.PathParam(c, "id")
	key := util.PathParam(c, "key")

	if len(key) == 0 {
		return api.Err(http.StatusBadRequest, "Invalid key", "The key must not be of length 0")
	}

	var data api.Metadata

	if err := util.ShouldBindJSONValidation(c, &data, false); err != nil {
		return api.Err(http.StatusBadRequest, "Invalid JSON", "%s", err)
	}

	if err := h.restream.SetProcessMetadata(id, key, data); err != nil {
		return api.Err(http.StatusNotFound, "Unknown process ID", "%s", err)
	}

	return c.JSON(http.StatusOK, data)
}

// GetMetadata returns the metadata stored with the Restreamer
// @Summary Retrieve JSON metadata from a key
// @Description Retrieve the previously stored JSON metadata under the given key. If the key is empty, all metadata will be returned.
// @Tags v16.7.2
// @ID metadata-3-get
// @Produce json
// @Param key path string true "Key for data store"
// @Success 200 {object} api.Metadata
// @Failure 404 {object} api.Error
// @Failure 400 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/metadata/{key} [get]
func (h *RestreamHandler) GetMetadata(c echo.Context) error {
	key := util.PathParam(c, "key")

	data, err := h.restream.GetMetadata(key)
	if err != nil {
		return api.Err(http.StatusNotFound, "Metadata not found", "%s", err)
	}

	return c.JSON(http.StatusOK, data)
}

// SetMetadata stores metadata with the Restreamer
// @Summary Add JSON metadata under the given key
// @Description Add arbitrary JSON metadata under the given key. If the key exists, all already stored metadata with this key will be overwritten. If the key doesn't exist, it will be created.
// @Tags v16.7.2
// @ID metadata-3-set
// @Produce json
// @Param key path string true "Key for data store"
// @Param data body api.Metadata true "Arbitrary JSON data"
// @Success 200 {object} api.Metadata
// @Failure 400 {object} api.Error
// @Security ApiKeyAuth
// @Router /api/v3/metadata/{key} [put]
func (h *RestreamHandler) SetMetadata(c echo.Context) error {
	key := util.PathParam(c, "key")

	if len(key) == 0 {
		return api.Err(http.StatusBadRequest, "Invalid key", "The key must not be of length 0")
	}

	var data api.Metadata

	if err := util.ShouldBindJSONValidation(c, &data, false); err != nil {
		return api.Err(http.StatusBadRequest, "Invalid JSON", "%s", err)
	}

	if err := h.restream.SetMetadata(key, data); err != nil {
		return api.Err(http.StatusBadRequest, "Invalid metadata", "%s", err)
	}

	return c.JSON(http.StatusOK, data)
}

func (h *RestreamHandler) getProcess(id, filterString string) (api.Process, error) {
	filter := strings.FieldsFunc(filterString, func(r rune) bool {
		return r == rune(',')
	})

	wants := map[string]bool{
		"config":   true,
		"state":    true,
		"report":   true,
		"metadata": true,
	}

	if len(filter) != 0 {
		for k := range wants {
			wants[k] = false
		}

		for _, f := range filter {
			if _, ok := wants[f]; ok {
				wants[f] = true
			}
		}
	}

	process, err := h.restream.GetProcess(id)
	if err != nil {
		return api.Process{}, err
	}

	info := api.Process{
		ID:        process.ID,
		Reference: process.Reference,
		Type:      "ffmpeg",
		CreatedAt: process.CreatedAt,
		UpdatedAt: process.UpdatedAt,
	}

	if wants["config"] {
		info.Config = &api.ProcessConfig{}
		info.Config.Unmarshal(process.Config)
	}

	if wants["state"] {
		state, err := h.restream.GetProcessState(id)
		if err != nil {
			return api.Process{}, err
		}

		info.State = &api.ProcessState{}
		info.State.Unmarshal(state)

	}

	if wants["report"] {
		log, err := h.restream.GetProcessLog(id)
		if err != nil {
			return api.Process{}, err
		}

		info.Report = &api.ProcessReport{}
		info.Report.Unmarshal(log)
	}

	if wants["metadata"] {
		data, err := h.restream.GetProcessMetadata(id, "")
		if err != nil {
			return api.Process{}, err
		}

		info.Metadata = api.NewMetadata(data)

	}

	return info, nil
}
