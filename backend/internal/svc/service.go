package svc

import (
	"handfree-work/octo-backup/internal/base/web_"
	"handfree-work/octo-backup/internal/modules/plugin"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Db      *gorm.DB
	Auth    web_.Config
	Plugins *plugin.Registry
}
