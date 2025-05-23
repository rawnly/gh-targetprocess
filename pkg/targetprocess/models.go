package targetprocess

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/rawnly/gh-targetprocess/templates"
)

type Assignable struct {
	ID          int     `json:"Id"`
	Name        string  `json:"Name"`
	Description *string `json:"Description,omitempty"`
}

func (a *Assignable) URL(baseURL string) string {
	return fmt.Sprintf("%s/entity/%d", baseURL, a.ID)
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
		URL:             a.URL(baseURL),
		Description:     description,
		DescriptionRows: rows,
	}

	if err = tmpl.Execute(&buf, payload); err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}

	return buf.String()
}
