package templates

import _ "embed"

//go:embed pr-body.tmpl
var PrBodyTemplate string

//go:embed pr-title.tmpl
var PrTitleTemplate string

// PRBodyTemplate is the embedded template string
func PRBodyTemplate() string {
	return PrBodyTemplate
}

// PRTitleTemplate is the embedded template string
func PRTitleTemplate() string {
	return PrTitleTemplate
}
