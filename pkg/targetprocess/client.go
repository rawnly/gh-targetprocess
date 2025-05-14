package targetprocess

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/rawnly/gh-targetprocess/templates"
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

type Assignable struct {
	ID          int     `json:"Id"`
	Name        string  `json:"Name"`
	Description *string `json:"Description,omitempty"`
}

func (a *Assignable) URL() string {
	return fmt.Sprintf("https://www.targetprocess.com/entity/%d", a.ID)
}

func (a *Assignable) GetPRTitle() string {
	return fmt.Sprintf("[%d]: %s", a.ID, a.Name)
}

func (a *Assignable) GetPRBody(baseURL string) string {
	var buf bytes.Buffer

	tmpl, err := template.New("pr-body").Parse(templates.PRBodyTemplate())
	if err != nil {
		return err.Error()
	}

	description := a.Description

	if description != nil {
		if strings.Contains(*description, "<!--markdown-->") {
			*description = strings.Replace(*description, "<!--markdown-->", "", 1)
		} else {
			var md string
			md, err = htmltomarkdown.ConvertString(*description)
			if err != nil {
				return err.Error()
			}

			*description = md
		}
	}

	var rows []string
	if description != nil {
		rows = strings.Split(*description, "\n\n")
	}

	payload := struct {
		ID              int
		Name            string
		URL             string
		Description     *string
		DescriptionRows []string
	}{
		ID:              a.ID,
		Name:            a.Name,
		URL:             a.URL(),
		Description:     description,
		DescriptionRows: rows,
	}

	if err = tmpl.Execute(&buf, payload); err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}

	return buf.String()
}

func New(baseURL string, apiKey string) *Client {
	return &Client{
		client:  newHTTPClient(),
		baseURL: fmt.Sprintf("%s/api", baseURL),
		apiKey:  apiKey,
	}
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
		return fmt.Errorf("error: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return err
	}

	return nil
}
