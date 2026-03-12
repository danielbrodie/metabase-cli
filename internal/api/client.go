package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/danielbrodie/metabase-cli/internal/config"
)

// APIError represents a non-2xx response from Metabase.
type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.Status, e.Message)
}

// Client is a thin Metabase API wrapper.
type Client struct {
	Profile *config.Profile
}

func New(profile *config.Profile) *Client {
	return &Client{Profile: profile}
}

func (c *Client) Do(method, path string, body interface{}, params map[string]string) (json.RawMessage, error) {
	endpoint := c.Profile.URL + "/api" + path

	if len(params) > 0 {
		v := url.Values{}
		for k, val := range params {
			if val != "" {
				v.Set(k, val)
			}
		}
		if len(v) > 0 {
			endpoint += "?" + v.Encode()
		}
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, endpoint, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.Profile.Token != "" {
		req.Header.Set("X-Metabase-Session", c.Profile.Token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		msg := string(respBody)
		var apiErr map[string]interface{}
		if json.Unmarshal(respBody, &apiErr) == nil {
			if m, ok := apiErr["message"].(string); ok && m != "" {
				msg = m
			}
		}
		return nil, &APIError{Status: resp.StatusCode, Message: msg}
	}

	return json.RawMessage(respBody), nil
}

func (c *Client) Get(path string, params map[string]string) (json.RawMessage, error) {
	return c.Do("GET", path, nil, params)
}

func (c *Client) Post(path string, body interface{}) (json.RawMessage, error) {
	return c.Do("POST", path, body, nil)
}

func (c *Client) Put(path string, body interface{}) (json.RawMessage, error) {
	return c.Do("PUT", path, body, nil)
}

func (c *Client) Delete(path string) (json.RawMessage, error) {
	return c.Do("DELETE", path, nil, nil)
}
