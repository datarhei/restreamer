package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/datarhei/core/v16/http/api"
	"github.com/datarhei/core/v16/http/mock"
	"github.com/stretchr/testify/require"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    interface{}
}

func getDummyRestreamHandler() (*RestreamHandler, error) {
	rs, err := mock.DummyRestreamer("../../mock")
	if err != nil {
		return nil, err
	}

	handler := NewRestream(rs)

	return handler, nil
}

func getDummyRestreamRouter() (*echo.Echo, error) {
	router := mock.DummyEcho()

	restream, err := getDummyRestreamHandler()
	if err != nil {
		return nil, err
	}

	router.GET("/", restream.GetAll)
	router.POST("/", restream.Add)
	router.GET("/:id", restream.Get)
	router.GET("/:id/report", restream.GetReport)
	router.PUT("/:id", restream.Update)
	router.DELETE("/:id", restream.Delete)
	router.PUT("/:id/command", restream.Command)

	return router, nil
}

func TestAddProcessMissingField(t *testing.T) {
	router, err := getDummyRestreamRouter()
	require.NoError(t, err)

	data := mock.Read(t, "./fixtures/addProcessMissingField.json")

	mock.Request(t, http.StatusOK, router, "POST", "/", data)
}

func TestAddProcessInvalidType(t *testing.T) {
	router, err := getDummyRestreamRouter()
	require.NoError(t, err)

	data := mock.Read(t, "./fixtures/addProcessInvalidType.json")

	mock.Request(t, http.StatusBadRequest, router, "POST", "/", data)
}

func TestAddProcess(t *testing.T) {
	router, err := getDummyRestreamRouter()
	require.NoError(t, err)

	data := mock.Read(t, "./fixtures/addProcess.json")

	response := mock.Request(t, http.StatusOK, router, "POST", "/", data)

	mock.Validate(t, &api.ProcessConfig{}, response.Data)
}

func TestUpdateProcessInvalid(t *testing.T) {
	router, err := getDummyRestreamRouter()
	require.NoError(t, err)

	data := mock.Read(t, "./fixtures/addProcess.json")

	response := mock.Request(t, http.StatusOK, router, "POST", "/", data)

	mock.Validate(t, &api.ProcessConfig{}, response.Data)

	update := bytes.Buffer{}
	_, err = update.ReadFrom(mock.Read(t, "./fixtures/addProcess.json"))
	require.NoError(t, err)

	proc := api.ProcessConfig{}
	err = json.Unmarshal(update.Bytes(), &proc)
	require.NoError(t, err)

	// invalid address
	proc.Output[0].Address = ""

	encoded, err := json.Marshal(&proc)
	require.NoError(t, err)

	update.Reset()
	_, err = update.Write(encoded)
	require.NoError(t, err)

	mock.Request(t, http.StatusBadRequest, router, "PUT", "/"+proc.ID, &update)
	mock.Request(t, http.StatusOK, router, "GET", "/"+proc.ID, nil)
}

func TestUpdateReplaceProcess(t *testing.T) {
	router, err := getDummyRestreamRouter()
	require.NoError(t, err)

	data := mock.Read(t, "./fixtures/addProcess.json")

	response := mock.Request(t, http.StatusOK, router, "POST", "/", data)

	mock.Validate(t, &api.ProcessConfig{}, response.Data)

	update := bytes.Buffer{}
	_, err = update.ReadFrom(mock.Read(t, "./fixtures/addProcess.json"))
	require.NoError(t, err)

	proc := api.ProcessConfig{}
	err = json.Unmarshal(update.Bytes(), &proc)
	require.NoError(t, err)

	encoded, err := json.Marshal(&proc)
	require.NoError(t, err)

	update.Reset()
	_, err = update.Write(encoded)
	require.NoError(t, err)

	response = mock.Request(t, http.StatusOK, router, "PUT", "/test", &update)

	mock.Validate(t, &api.ProcessConfig{}, response.Data)

	mock.Request(t, http.StatusOK, router, "GET", "/test", nil)
}

func TestUpdateNewProcess(t *testing.T) {
	router, err := getDummyRestreamRouter()
	require.NoError(t, err)

	data := mock.Read(t, "./fixtures/addProcess.json")

	response := mock.Request(t, http.StatusOK, router, "POST", "/", data)

	mock.Validate(t, &api.ProcessConfig{}, response.Data)

	update := bytes.Buffer{}
	_, err = update.ReadFrom(mock.Read(t, "./fixtures/addProcess.json"))
	require.NoError(t, err)

	proc := api.ProcessConfig{}
	err = json.Unmarshal(update.Bytes(), &proc)
	require.NoError(t, err)

	proc.ID = "test2"

	encoded, err := json.Marshal(&proc)
	require.NoError(t, err)

	update.Reset()
	_, err = update.Write(encoded)
	require.NoError(t, err)

	response = mock.Request(t, http.StatusOK, router, "PUT", "/test", &update)

	mock.Validate(t, &api.ProcessConfig{}, response.Data)

	mock.Request(t, http.StatusNotFound, router, "GET", "/test", nil)
	mock.Request(t, http.StatusOK, router, "GET", "/test2", nil)
}

func TestUpdateNonExistentProcess(t *testing.T) {
	router, err := getDummyRestreamRouter()
	require.NoError(t, err)

	data := mock.Read(t, "./fixtures/addProcess.json")

	mock.Request(t, http.StatusNotFound, router, "PUT", "/test", data)
}

func TestRemoveUnknownProcess(t *testing.T) {
	router, err := getDummyRestreamRouter()
	require.NoError(t, err)

	mock.Request(t, http.StatusNotFound, router, "DELETE", "/foobar", nil)
}

func TestRemoveProcess(t *testing.T) {
	router, err := getDummyRestreamRouter()
	require.NoError(t, err)

	data := mock.Read(t, "./fixtures/removeProcess.json")

	mock.Request(t, http.StatusOK, router, "POST", "/", data)
	mock.Request(t, http.StatusOK, router, "DELETE", "/test", nil)
}

func TestProcessInfo(t *testing.T) {
	router, err := getDummyRestreamRouter()
	require.NoError(t, err)

	data := mock.Read(t, "./fixtures/addProcess.json")

	mock.Request(t, http.StatusOK, router, "POST", "/", data)
	response := mock.Request(t, http.StatusOK, router, "GET", "/test", nil)

	mock.Validate(t, &api.Process{}, response.Data)
}

func TestProcessReportNotFound(t *testing.T) {
	router, err := getDummyRestreamRouter()
	require.NoError(t, err)

	mock.Request(t, http.StatusNotFound, router, "GET", "/test/report", nil)
}

func TestProcessReport(t *testing.T) {
	router, err := getDummyRestreamRouter()
	require.NoError(t, err)

	data := mock.Read(t, "./fixtures/addProcess.json")

	mock.Request(t, http.StatusOK, router, "POST", "/", data)
	response := mock.Request(t, http.StatusOK, router, "GET", "/test/report", nil)

	mock.Validate(t, &api.ProcessReport{}, response.Data)
}

func TestProcessCommandNotFound(t *testing.T) {
	router, err := getDummyRestreamRouter()
	require.NoError(t, err)

	command := mock.Read(t, "./fixtures/commandStart.json")
	mock.Request(t, http.StatusBadRequest, router, "PUT", "/test/command", command)
}

func TestProcessCommandInvalid(t *testing.T) {
	router, err := getDummyRestreamRouter()
	require.NoError(t, err)

	data := mock.Read(t, "./fixtures/addProcess.json")

	mock.Request(t, http.StatusOK, router, "POST", "/", data)

	command := mock.Read(t, "./fixtures/commandInvalid.json")
	mock.Request(t, http.StatusBadRequest, router, "PUT", "/test/command", command)
}

func TestProcessCommand(t *testing.T) {
	router, err := getDummyRestreamRouter()
	require.NoError(t, err)

	data := mock.Read(t, "./fixtures/addProcess.json")

	mock.Request(t, http.StatusOK, router, "POST", "/", data)

	command := mock.Read(t, "./fixtures/commandStart.json")
	mock.Request(t, http.StatusOK, router, "PUT", "/test/command", command)
	mock.Request(t, http.StatusOK, router, "GET", "/test", data)

	command = mock.Read(t, "./fixtures/commandStop.json")
	mock.Request(t, http.StatusOK, router, "PUT", "/test/command", command)
	mock.Request(t, http.StatusOK, router, "GET", "/test", data)
}
