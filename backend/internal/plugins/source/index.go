package source

import (
	"fmt"

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
			return fmt.Errorf("加载备份来源插件: %w", err)
		}
		if err := r.Register(d); err != nil {
			return err
		}
	}
	return nil
}
