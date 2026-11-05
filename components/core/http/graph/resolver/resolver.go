package resolver

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/datarhei/core/v16/http/graph/models"
	"github.com/datarhei/core/v16/log"
	"github.com/datarhei/core/v16/monitor"
	"github.com/datarhei/core/v16/restream"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Restream  restream.Restreamer
	Monitor   monitor.HistoryReader
	LogBuffer log.BufferWriter
}

func (r *queryResolver) getProcess(id string) (*models.Process, error) {
	process, err := r.Restream.GetProcess(id)
	if err != nil {
		return nil, err
	}

	state, err := r.Restream.GetProcessState(id)
	if err != nil {
		return nil, err
	}

	report, err := r.Restream.GetProcessLog(id)
	if err != nil {
		return nil, err
	}

	m, err := r.Restream.GetProcessMetadata(id, "")
	if err != nil {
		return nil, err
	}

	metadata, ok := m.(map[string]interface{})
	if !ok {
		metadata = nil
	}

	p := &models.Process{}
	p.UnmarshalRestream(process, state, report, metadata)

	return p, nil
}

func (r *queryResolver) playoutRequest(method, addr, path, contentType string, data []byte) ([]byte, error) {
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

	defer response.Body.Close()

	// Read the whole response
	data, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
