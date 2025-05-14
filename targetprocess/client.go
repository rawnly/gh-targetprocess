package targetprocess

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"
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

	tmpl, err := template.New("pr-body.tmpl").Parse(`
See: [[{{.ID}}] {{.Name}}]({{.URL}})

{{ if .Description }}
{{ .Description }}
{{ end }}
		`)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	description := a.Description

	// Remove the <!--markdown--> comment from the description and cleanup
	if description != nil {
		tempDescription := strings.Replace(*description, "<!--markdown-->", "", 1)
		rows := strings.Split(tempDescription, "\n\n")
		for i, row := range rows {
			rows[i] = fmt.Sprintf("> %s", row)
		}

		tempDescription = strings.Join(rows, "\n>\n")

		description = &tempDescription
	}

	payload := struct {
		ID          int
		Name        string
		URL         string
		Description *string
	}{
		ID:          a.ID,
		Name:        a.Name,
		URL:         a.URL(),
		Description: description,
	}

	if err = tmpl.Execute(&buf, payload); err != nil {
		fmt.Println(err.Error())
		return ""
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
