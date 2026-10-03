package notification

import (
	"fmt"
	"handfree-work/octo-backup/internal/modules/plugin"
	"handfree-work/octo-backup/internal/plugins/notification/dingtalk"
	"handfree-work/octo-backup/internal/plugins/notification/email"
	"handfree-work/octo-backup/internal/plugins/notification/feishu"
	"handfree-work/octo-backup/internal/plugins/notification/webhook"
	"handfree-work/octo-backup/internal/plugins/notification/wecom"
)

type provider interface {
	Definition() (*plugin.Definition, error)
}

func Register(r *plugin.Registry) error {
	providers := []provider{email.Provider{}, dingtalk.Provider{}, feishu.Provider{}, wecom.Provider{}, webhook.Provider{}}
	for _, p := range providers {
		d, err := p.Definition()
		if err != nil {
			return fmt.Errorf("加载通知插件: %w", err)
		}
		if err := r.Register(d); err != nil {
			return err
		}
	}
	return nil
}
