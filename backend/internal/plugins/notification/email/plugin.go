package email

import (
	"context"
	_ "embed"
	"handfree-work/octo-backup/internal/modules/plugin"
)

//go:embed metadata.yaml
var metadata []byte

type Provider struct{}

func (Provider) Definition() (*plugin.Definition, error) {
	return plugin.NewDefinition(metadata, &EmailNotifier{})
}

type EmailNotifier struct {
	PluginContext *plugin.PluginContext
	Host          string
	Port          int
	Username      string
	Password      string
	From          string
}

func (EmailNotifier) ExecuteAction(_ context.Context, action string, _ map[string]any) (any, error) {
	if action == "onSend" {
		return map[string]any{"sent": true}, nil
	}
	return nil, plugin.NewActionNotFoundError()
}
