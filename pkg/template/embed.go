package template

import "embed"

//go:embed *.go.tmpl predicate/*.go.tmpl __entity-package__/*.go.tmpl
var files embed.FS
