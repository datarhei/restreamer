package zerossl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client acts as a ZeroSSL API client. It facilitates ZeroSSL certificate operations.
type Client struct {
	// REQUIRED: Your ZeroSSL account access key.
	AccessKey string `json:"access_key"`

	// Optionally adjust the base URL of the API.
	// Default: https://api.zerossl.com
	BaseURL string `json:"base_url,omitempty"`

	// Optionally configure a custom HTTP client.
	HTTPClient *http.Client `json:"-"`
}

func (c Client) httpGet(ctx context.Context, endpoint string, qs url.Values, target any) error {
	url := c.url(endpoint, qs)
	return c.httpRequest(ctx, http.MethodGet, url, nil, target)
}

func (c Client) httpPost(ctx context.Context, endpoint string, qs url.Values, payload, target any) error {
	var reqBody io.Reader
	if payload != nil {
		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(payloadJSON)
	}
	url := c.url(endpoint, qs)
	return c.httpRequest(ctx, http.MethodPost, url, reqBody, target)
}

func (c Client) httpRequest(ctx context.Context, method, reqURL string, reqBody io.Reader, target any) error {
	r, err := http.NewRequestWithContext(ctx, method, reqURL, reqBody)
	if err != nil {
		return err
	}
	if reqBody != nil {
		r.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient().Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// because the ZeroSSL API doesn't use HTTP status codes to indicate an error,
	// nor does each response body have a consistent way of detecting success/error,
	// we have to implement a hack: download the entire response body and try
	// decoding it as JSON in a way that errors if there's any unknown fields
	// (such as "success"), because if there is an unkown field, either our model
	// is outdated, or there was an error payload in the response instead of the
	// expected structure, so we then try again to decode to an error struct
	respBytes, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024*2))
	if err != nil {
		return fmt.Errorf("failed reading response body: %v", err)
	}

	// assume success first by trying to decode payload into output target
	dec := json.NewDecoder(bytes.NewReader(respBytes))
	dec.DisallowUnknownFields() // important hacky hack so we can detect an error payload
	originalDecodeErr := dec.Decode(&target)
	if originalDecodeErr == nil {
		return nil
	}

	// could have gotten any kind of error, really; but assuming valid JSON,
	// most likely it is an error payload
	var apiError APIError
	if err := json.NewDecoder(bytes.NewReader(respBytes)).Decode(&apiError); err != nil {
		return fmt.Errorf("request succeeded, but decoding JSON response failed: %v (raw=%s)", err, respBytes)
	}

	// successfully got an error! or did we?
	if apiError.Success {
		return apiError // ummm... why are we getting an error if it was successful ??? is this not really an error?
	}

	// remove access_key from URL so it doesn't leak into logs
	u, err := url.Parse(reqURL)
	if err != nil {
		reqURL = fmt.Sprintf("<invalid url: %v>", err)
	}
	if u != nil {
		q, err := url.ParseQuery(u.RawQuery)
		if err == nil {
			q.Set(accessKeyParam, "redacted")
			u.RawQuery = q.Encode()
			reqURL = u.String()
		}
	}

	return fmt.Errorf("%s %s: HTTP %d: %v (raw=%s decode_error=%v)", method, reqURL, resp.StatusCode, apiError, respBytes, originalDecodeErr)
}

func (c Client) url(endpoint string, qs url.Values) string {
	baseURL := c.BaseURL
	if baseURL == "" {
		baseURL = BaseURL
	}

	// for consistency, ensure endpoint starts with /
	// and base URL does NOT end with /.
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}
	baseURL = strings.TrimSuffix(baseURL, "/")

	if qs == nil {
		qs = url.Values{}
	}
	qs.Set(accessKeyParam, c.AccessKey)

	return fmt.Sprintf("%s%s?%s", baseURL, endpoint, qs.Encode())
}

func (c Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return httpClient
}

var httpClient = &http.Client{
	Timeout: 2 * time.Minute,
}

// anyBool is a hacky type that accepts true or 1 (or their string variants),
// or "yes" or "y", and any casing variants of the same, as a boolean true when
// unmarshaling JSON. Everything else is boolean false.
//
// This is needed due to type inconsistencies in ZeroSSL's API with "success" values.
type anyBool bool

// UnmarshalJSON satisfies json.Unmarshaler according to
// this type's documentation.
func (ab *anyBool) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return io.EOF
	}
	switch strings.ToLower(string(b)) {
	case `true`, `"true"`, `1`, `"1"`, `"yes"`, `"y"`:
		*ab = true
	}
	return nil
}

// MarshalJSON marshals ab to either true or false.
func (ab *anyBool) MarshalJSON() ([]byte, error) {
	if ab != nil && *ab {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}

const accessKeyParam = "access_key"
