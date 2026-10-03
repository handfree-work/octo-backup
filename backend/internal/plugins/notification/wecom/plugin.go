package wecom

import (
	_ "embed"
	"handfree-work/octo-backup/internal/modules/plugin"
)

//go:embed metadata.yaml
var metadata []byte

type Provider struct{}

func (Provider) Definition() (*plugin.Definition, error) {
	return plugin.NewGenericDefinition(metadata)
}
