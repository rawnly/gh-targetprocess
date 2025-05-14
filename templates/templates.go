package templates

import _ "embed"

//go:embed pr-body.tmpl
var PrBodyTemplate string

// PRBodyTemplate is the embedded template string
func PRBodyTemplate() string {
	return PrBodyTemplate
}
