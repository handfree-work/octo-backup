package minio

import (
	"context"
	_ "embed"
	"handfree-work/octo-backup/internal/modules/plugin"
)

//go:embed metadata.yaml
var metadata []byte

type Provider struct{}

func (Provider) Definition() (*plugin.Definition, error) {
	return plugin.NewDefinition(metadata, func(map[string]any) (plugin.ActionExecutor, error) { return executor{}, nil })
}
type executor struct{}
func (executor) ExecuteAction(_ context.Context, action string, _ map[string]any) (any, error) { switch action { case "onTest": return map[string]any{"ok": true}, nil; case "onListBuckets": return []any{}, nil }; return nil, plugin.ErrActionNotFound }
