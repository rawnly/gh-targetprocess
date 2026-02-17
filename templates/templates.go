package templates

import (
	"bytes"
	_ "embed"
	"os"

	"github.com/rawnly/gh-targetprocess/internal/utils"
)

//go:embed pr-body.tmpl
var PrBodyTemplate string

//go:embed pr-title.tmpl
var PrTitleTemplate string

// PRBodyTemplate is the embedded template string
func PRBodyTemplate() (string, error) {
	path := utils.ExpandPath("~/.config/gh-targetprocess/pr-body.tmpl")
	template, err := os.ReadFile(path)
	if err != nil {
		return PrBodyTemplate, nil
	}

	return string(template), nil
}

// PRTitleTemplate is the embedded template string
func PRTitleTemplate() (string, error) {
	path := utils.ExpandPath("~/.config/gh-targetprocess/pr-title.tmpl")
	template, err := os.ReadFile(path)
	if err != nil {
		return PrTitleTemplate, nil
	}

	return string(template), nil
}

func WriteDefaults() error {
	if err := os.MkdirAll(utils.ExpandPath("~/.config/gh-targetproces"), 0700); err != nil {
		return err
	}

	titleBuffer := bytes.NewBufferString(PrTitleTemplate)
	if err := os.WriteFile(utils.ExpandPath("~/.config/gh-targetprocess/pr-title.tmpl"), titleBuffer.Bytes(), 0600); err != nil {
		return err
	}

	bodyBuffer := bytes.NewBufferString(PrBodyTemplate)
	if err := os.WriteFile(utils.ExpandPath("~/.config/gh-targetprocess/pr-body.tmpl"), bodyBuffer.Bytes(), 0600); err != nil {
		return err
	}

	return nil
}
