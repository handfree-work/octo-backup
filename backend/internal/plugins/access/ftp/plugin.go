package ftp

import (
	"context"
	_ "embed"
	"handfree-work/octo-backup/internal/modules/plugin"
)

//go:embed metadata.yaml
var metadata []byte

type Provider struct{}

func (Provider) Definition() (*plugin.Definition, error) {
	return plugin.NewDefinition(metadata, &FtpAccess{})
}

type FtpAccess struct {
	PluginContext *plugin.PluginContext
	Host          string
	Port          int
	Username      string
	Password      string
}

func (FtpAccess) ExecuteAction(_ context.Context, action string, _ map[string]any) (any, error) {
	if action == "onTest" {
		return map[string]any{"ok": true}, nil
	}
	return nil, plugin.NewActionNotFoundError()
}
