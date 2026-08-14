package svc

import (
	"handfree-work/web-restic/internal/auth"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Db   *gorm.DB
	Auth auth.Config
}
