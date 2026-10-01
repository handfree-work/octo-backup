package models

import "handfree-work/web-restic/internal/base/db_"

// StorageRepository 描述一个 Restic 存储仓库。
type StorageRepository struct {
	db_.BaseModel
	Name        string `json:"name" gorm:"size:100;not null;uniqueIndex;comment:仓库名称"`
	Repository  string `json:"repository" gorm:"size:500;not null;comment:仓库地址"`
	Password    string `json:"-" gorm:"size:255;comment:仓库密码"`
	Description string `json:"description" gorm:"size:500;comment:备注"`
}
