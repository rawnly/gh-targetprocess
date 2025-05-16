package targetprocess

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type transport struct {
	headers map[string]string
	base    http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header.Add(k, v)
	}

	base := t.base

	if base == nil {
		base = http.DefaultTransport
	}

	return base.RoundTrip(req)
}

func newHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &transport{
			headers: map[string]string{
				"Accept": "application/json",
			},
		},
	}

	return client
}

type Client struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

func New(baseURL string, apiKey string) *Client {
	return &Client{
		client:  newHTTPClient(),
		baseURL: fmt.Sprintf("%s/api", baseURL),
		apiKey:  apiKey,
	}
}

func (c *Client) GetAssignable(id string) (*Assignable, error) {
	var assignable Assignable

	path := fmt.Sprintf("/v1/assignable/%s", id)

	if err := c.Get(path, &assignable); err != nil {
		return nil, err
	}

	return &assignable, nil
}

func (c *Client) Test(path string) error {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return err
	}

	params := req.URL.Query()
	params.Add("access_token", c.apiKey)

	req.URL.RawQuery = params.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	return nil
}

func (c *Client) Get(path string, response any) error {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return err
	}

	params := req.URL.Query()
	params.Add("access_token", c.apiKey)

	req.URL.RawQuery = params.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return err
	}

	return nil
}
