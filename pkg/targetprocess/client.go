package targetprocess

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
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
		Timeout: time.Second * 30,
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

func (c *Client) GetAssignable(ctx context.Context, id string) (*Assignable, error) {
	var assignable Assignable

	path := fmt.Sprintf("/v1/assignable/%s", id)

	if err := c.Get(ctx, path, &assignable); err != nil {
		return nil, err
	}

	return &assignable, nil
}

func (c *Client) UpdateState(ctx context.Context, assignableId int, newState EntityState) error {
	path := fmt.Sprintf("/v1/assignable/%d", assignableId)

	payload := map[string]any{
		"EntityState": map[string]any{
			"Id": newState,
		},
	}

	return c.Post(ctx, path, payload)
}

func (c *Client) PostComment(ctx context.Context, content string, assignableId int) error {
	path := "/v1/comments"

	payload := CreateCommentPayload{
		Description: content,
		General: struct {
			Id int "json:\"Id\""
		}{
			Id: assignableId,
		},
	}

	return c.Post(ctx, path, payload)
}

func (c *Client) Test(ctx context.Context, path string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+path, nil)
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

func (c *Client) Get(ctx context.Context, path string, response any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+path, nil)
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

func (c *Client) Post(ctx context.Context, path string, body any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	payload := bytes.NewBuffer(b)

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, payload)
	if err != nil {
		return err
	}

	// Query
	params := req.URL.Query()
	params.Add("access_token", c.apiKey)

	req.URL.RawQuery = params.Encode()

	// Headers
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
