package access

import (
	"fmt"
	"handfree-work/octo-backup/internal/modules/plugin"
	"handfree-work/octo-backup/internal/plugins/access/aliyun"
	"handfree-work/octo-backup/internal/plugins/access/ftp"
	"handfree-work/octo-backup/internal/plugins/access/minio"
	"handfree-work/octo-backup/internal/plugins/access/s3"
	"handfree-work/octo-backup/internal/plugins/access/sftp"
	"handfree-work/octo-backup/internal/plugins/access/ssh"
	"handfree-work/octo-backup/internal/plugins/access/tencent"
)

type provider interface {
	Definition() (*plugin.Definition, error)
}

func Register(r *plugin.Registry) error {
	providers := []provider{ssh.Provider{}, ftp.Provider{}, sftp.Provider{}, aliyun.Provider{}, tencent.Provider{}, s3.Provider{}, minio.Provider{}}
	for _, p := range providers {
		d, err := p.Definition()
		if err != nil {
			return fmt.Errorf("加载授权插件: %w", err)
		}
		if err := r.Register(d); err != nil {
			return err
		}
	}
	return nil
}
