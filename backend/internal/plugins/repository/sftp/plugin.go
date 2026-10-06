package sftp

import (
	"context"
	_ "embed"
	"handfree-work/octo-backup/internal/modules/plugin"
)

//go:embed metadata.yaml
var metadata []byte

type Provider struct{}

func (Provider) Definition() (*plugin.Definition, error) {
	return plugin.NewDefinition(metadata, &SftpRepository{})
}

type SftpRepository struct {
	PluginContext *plugin.PluginContext
	Path          string
	AccessId      int
	Password      string
}

func (SftpRepository) ExecuteAction(_ context.Context, action string, _ map[string]any) (any, error) {
	if action == "onBuild" {
		return map[string]any{"ok": true}, nil
	}
	return nil, plugin.NewActionNotFoundError()
}
