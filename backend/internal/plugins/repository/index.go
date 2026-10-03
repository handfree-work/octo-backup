package repository

import (
	"fmt"
	"handfree-work/octo-backup/internal/modules/plugin"
	"handfree-work/octo-backup/internal/plugins/repository/local"
	"handfree-work/octo-backup/internal/plugins/repository/minio"
	"handfree-work/octo-backup/internal/plugins/repository/s3"
	"handfree-work/octo-backup/internal/plugins/repository/sftp"
)

type provider interface {
	Definition() (*plugin.Definition, error)
}

func Register(r *plugin.Registry) error {
	providers := []provider{local.Provider{}, sftp.Provider{}, s3.Provider{}, minio.Provider{}}
	for _, p := range providers {
		d, err := p.Definition()
		if err != nil {
			return fmt.Errorf("加载仓库插件: %w", err)
		}
		if err := r.Register(d); err != nil {
			return err
		}
	}
	return nil
}
