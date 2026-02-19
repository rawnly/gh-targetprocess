package templates

import (
	"bytes"
	_ "embed"
	"os"

	"github.com/rawnly/gh-targetprocess/internal/utils"
)

const bodyTemplateFileName = "pr-body.tmpl"
const titleTemplateFileName = "pr-title.tmpl"

//go:embed pr-body.tmpl
var prBodyDefaultTemplate string

//go:embed pr-title.tmpl
var prTitleDefaultTemplate string

// PRBodyTemplate is the embedded template string
func PRBodyTemplate() (string, error) {
	path := utils.GetConfigFilePath(bodyTemplateFileName)
	template, err := os.ReadFile(path)
	if err != nil {
		return prBodyDefaultTemplate, nil
	}

	return string(template), nil
}

// PRTitleTemplate is the embedded template string
func PRTitleTemplate() (string, error) {
	path := utils.GetConfigFilePath(titleTemplateFileName)
	template, err := os.ReadFile(path)
	if err != nil {
		return prTitleDefaultTemplate, nil
	}

	return string(template), nil
}

func WriteDefaults() error {
	if err := os.MkdirAll(utils.ExpandPath(utils.ConfigDir), 0700); err != nil {
		return err
	}

	titleBuffer := bytes.NewBufferString(prTitleDefaultTemplate)
	if err := os.WriteFile(utils.GetConfigFilePath(titleTemplateFileName), titleBuffer.Bytes(), 0600); err != nil {
		return err
	}

	bodyBuffer := bytes.NewBufferString(prBodyDefaultTemplate)
	if err := os.WriteFile(utils.GetConfigFilePath(bodyTemplateFileName), bodyBuffer.Bytes(), 0600); err != nil {
		return err
	}

	return nil
}
