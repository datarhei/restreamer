package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/datarhei/core/v16/log"
)

type API interface {
	Monitor(id string, data MonitorData) (MonitorResponse, error)
}

type Config struct {
	URL    string
	Token  string
	Client *http.Client
	Logger log.Logger
}

type api struct {
	url   string
	token string

	accessToken     string
	accessTokenType string

	client *http.Client

	logger log.Logger
}

func New(config Config) (API, error) {
	a := &api{
		url:    config.URL,
		token:  config.Token,
		client: config.Client,
		logger: config.Logger,
	}

	if a.logger == nil {
		a.logger = log.New("")
	}

	if !strings.HasSuffix(a.url, "/") {
		a.url = a.url + "/"
	}

	if a.client == nil {
		a.client = &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 30 * time.Second,
			},
			Timeout: 5 * time.Second,
		}
	}

	return a, nil
}

var (
	errStatusAuthorization = errors.New("authorization with token failed")
)

type statusError struct {
	code     int
	response []byte
}

func (e statusError) Error() string {
	return fmt.Sprintf("response code %d (%s)", e.code, e.response)
}

func (e statusError) Is(target error) bool {
	_, ok := target.(statusError)

	return ok
}

type copyReader struct {
	reader io.Reader
	copy   *bytes.Buffer
}

func newCopyReader(r io.Reader) io.Reader {
	c := &copyReader{
		reader: r,
		copy:   new(bytes.Buffer),
	}

	return c
}

func (c *copyReader) Read(p []byte) (int, error) {
	i, err := c.reader.Read(p)

	c.copy.Write(p)

	if err == io.EOF {
		c.reader = c.copy
		c.copy = &bytes.Buffer{}
	}

	return i, err
}

func (a *api) callWithRetry(method, path string, body io.Reader) ([]byte, error) {
	if len(a.accessToken) == 0 {
		err := a.refreshAccessToken(a.token)
		if err != nil {
			return nil, fmt.Errorf("invalid token: %w", err)
		}
	}

	cpr := newCopyReader(body)

	data, err := a.call(method, path, cpr)
	if errors.Is(err, errStatusAuthorization) {
		err = a.refreshAccessToken(a.token)
		if err != nil {
			return nil, fmt.Errorf("invalid token: %w", err)
		}

		data, err = a.call(method, path, cpr)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return data, err
}

func (a *api) refreshAccessToken(refreshToken string) error {
	req, err := http.NewRequest(http.MethodPut, a.url+"api/v1/token/login", nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+refreshToken)

	res, err := a.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		if res.StatusCode == 401 {

			return errStatusAuthorization
		}

		data, _ := io.ReadAll(res.Body)
		return statusError{
			code:     res.StatusCode,
			response: data,
		}
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return statusError{
			code: res.StatusCode,
		}
	}

	token := tokenResponse{}

	if err := json.Unmarshal(data, &token); err != nil {
		return fmt.Errorf("error parsing response: %w", err)
	}

	a.accessToken = token.AccessToken
	a.accessTokenType = token.TokenType

	return nil
}

func (a *api) call(method, path string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, a.url+path, body)
	if err != nil {
		return nil, err
	}

	if len(a.accessToken) != 0 {
		req.Header.Add("Authorization", a.accessTokenType+" "+a.accessToken)
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		if res.StatusCode == 401 {
			return nil, errStatusAuthorization
		}

		data, _ := io.ReadAll(res.Body)
		return nil, statusError{
			code:     res.StatusCode,
			response: data,
		}
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", statusError{
			code: res.StatusCode,
		})
	}

	return data, nil
}

func (a *api) Monitor(id string, monitordata MonitorData) (MonitorResponse, error) {
	var data bytes.Buffer
	encoder := json.NewEncoder(&data)
	if err := encoder.Encode(monitordata); err != nil {
		return MonitorResponse{}, err
	}

	/*
		b, err := json.MarshalIndent(monitordata, "", "  ")
		if err == nil {
			fmt.Println(string(b))
		}
	*/

	response, err := a.callWithRetry(http.MethodPut, "api/v1/core/monitor/"+id, &data)
	if err != nil {
		return MonitorResponse{}, fmt.Errorf("error sending request: %w", err)
	}

	r := MonitorResponse{}

	if err := json.Unmarshal(response, &r); err != nil {
		return MonitorResponse{}, fmt.Errorf("error parsing response: %w", err)
	}

	return r, nil
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type MonitorResponse struct {
	Next uint64 `json:"next_update"`
}

type MonitorProcessData struct {
	ID     string              `json:"id"`
	RefID  string              `json:"id_ref"`
	CPU    []json.Number       `json:"cpu"`
	Mem    []json.Number       `json:"mem"`
	Uptime uint64              `json:"uptime_sec"`
	Output map[string][]uint64 `json:"output"`
}

type MonitorData struct {
	Version       string                `json:"version"`
	Uptime        uint64                `json:"uptime_sec"` // seconds
	SysCPU        []json.Number         `json:"sys_cpu"`
	SysMemory     []json.Number         `json:"sys_mem"`  // bytes
	SysDisk       []json.Number         `json:"sys_disk"` // bytes
	FSMem         []json.Number         `json:"fs_mem"`   // bytes
	FSDisk        []json.Number         `json:"fs_disk"`  // bytes
	NetTX         []json.Number         `json:"net_tx"`   // kbit/s
	Session       []json.Number         `json:"viewer"`
	ProcessStates [6]uint64             `json:"proc_states"` // finished, starting, running, finishing, failed, killed
	Processes     *[]MonitorProcessData `json:"procs,omitempty"`
}
