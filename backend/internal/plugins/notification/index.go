package notification

import (
	"handfree-work/octo-backup/internal/base/error_"
	"handfree-work/octo-backup/internal/modules/plugin"
	"handfree-work/octo-backup/internal/plugins/notification/email"
)

type provider interface {
	Definition() (*plugin.Definition, error)
}

func Register(r *plugin.Registry) error {
	providers := []provider{email.Provider{}}
	for _, p := range providers {
		d, err := p.Definition()
		if err != nil {
			return error_.NewWrapError("加载通知插件", err)
		}
		if err := r.Register(d); err != nil {
			return err
		}
	}
	return nil
}
