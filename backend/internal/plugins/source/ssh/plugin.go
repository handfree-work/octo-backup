package ssh

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"handfree-work/octo-backup/internal/base/error_"
	"handfree-work/octo-backup/internal/modules/plugin"
	sshaccess "handfree-work/octo-backup/internal/plugins/access/ssh"
)

//go:embed metadata.yaml
var metadata []byte

type Provider struct{}

func (Provider) Definition() (*plugin.Definition, error) {
	return plugin.NewDefinition(metadata, func(config map[string]any) (plugin.ActionExecutor, error) {
		return &executor{config: config}, nil
	})
}

type executor struct {
	config map[string]any
}

func (e *executor) ExecuteAction(ctx context.Context, action string, params map[string]any) (any, error) {
	form := make(map[string]any, len(e.config)+len(params))
	for key, value := range e.config {
		form[key] = value
	}
	for key, value := range params {
		form[key] = value
	}
	switch action {
	case "onListPaths":
		result, err := e.onListPaths(ctx, form)
		if err != nil {
			return nil, error_.NewTextError(fmt.Sprintf("source.ssh.onListPaths: %v", err))
		}
		return result, nil
	default:
		return nil, plugin.ErrActionNotFound
	}
}

// onListPaths 根据完整插件表单读取 SSH 主机目录。
// 目录遍历接入 SSH 执行器后，只需替换此方法，action 分发和表单参数保持不变。
func (e *executor) onListPaths(ctx context.Context, form map[string]any) (any, error) {
	if form["accessId"] == nil {
		return nil, error_.NewTextError("SSH 授权 ID 不能为空")
	}
	access, ok := form["access"].(map[string]any)
	if !ok {
		return nil, error_.NewTextError("SSH 授权配置不存在")
	}
	config, ok := access["config"].(map[string]any)
	if !ok {
		return nil, error_.NewTextError("SSH 授权配置无效")
	}
	client := sshaccess.NewSshAccess(config)
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}
	defer client.Close()
	paths, err := client.ListDirectories(strings.TrimSpace(fmt.Sprint(form["path"])))
	if err != nil {
		return nil, err
	}
	return map[string]any{"paths": paths}, nil
}
