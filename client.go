package hackeronecli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const defaultBaseURL = "https://api.hackerone.com/v1"

var Version = "dev"

type APIError struct {
	StatusCode int      `json:"status_code"`
	Message    string   `json:"message"`
	Errors     []string `json:"errors"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("API error %d", e.StatusCode)
}

type Client struct {
	BaseURL    string
	Identifier string
	Token      string
	HTTPClient *http.Client
}

func NewClient(identifier, token string) *Client {
	return &Client{
		BaseURL:    defaultBaseURL,
		Identifier: identifier,
		Token:      token,
		HTTPClient: http.DefaultClient,
	}
}

func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	u := c.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.Identifier, c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "h1-cli/"+Version)
	return req, nil
}

func (c *Client) Get(ctx context.Context, path string, params url.Values) (*http.Response, error) {
	if len(params) > 0 {
		path = path + "?" + params.Encode()
	}
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req)
}

func (c *Client) Post(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req)
}

func (c *Client) Put(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req)
}

func (c *Client) Patch(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req)
}

func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req)
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if decodeErr := json.NewDecoder(resp.Body).Decode(apiErr); decodeErr != nil {
			apiErr.Message = http.StatusText(resp.StatusCode)
		}
		return nil, apiErr
	}
	return resp, nil
}

func decodeResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}
