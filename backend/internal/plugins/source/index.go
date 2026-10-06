package source

import (
	"handfree-work/octo-backup/internal/base/error_"

	"handfree-work/octo-backup/internal/modules/plugin"
	"handfree-work/octo-backup/internal/plugins/source/ssh"
)

type provider interface {
	Definition() (*plugin.Definition, error)
}

func Register(r *plugin.Registry) error {
	for _, p := range []provider{ssh.Provider{}} {
		d, err := p.Definition()
		if err != nil {
			return error_.NewWrapError("加载备份来源插件", err)
		}
		if err := r.Register(d); err != nil {
			return err
		}
	}
	return nil
}
