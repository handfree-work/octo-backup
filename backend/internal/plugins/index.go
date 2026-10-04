package plugins

import (
	"handfree-work/octo-backup/internal/modules/plugin"
	"handfree-work/octo-backup/internal/plugins/access"
	"handfree-work/octo-backup/internal/plugins/notification"
	"handfree-work/octo-backup/internal/plugins/repository"
	"handfree-work/octo-backup/internal/plugins/source"
)

func RegisterAll(r *plugin.Registry) error {
	if err := access.Register(r); err != nil {
		return err
	}
	if err := repository.Register(r); err != nil {
		return err
	}
	if err := source.Register(r); err != nil {
		return err
	}
	return notification.Register(r)
}
