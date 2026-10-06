package models

import "handfree-work/octo-backup/internal/base/db_"

type BackupPlan struct {
	db_.BaseModel
	Name         string `json:"name" gorm:"size:100;not null"`
	SourceId     int64  `json:"sourceId" gorm:"not null;index"`
	RepositoryId int64  `json:"repositoryId" gorm:"not null;index:idx_backup_plan_repo_path,unique"`
	RepoTag      string `json:"repoTag" gorm:"size:1000;index:idx_backup_plan_repo_path,unique"`
	Schedule     string `json:"schedule" gorm:"size:100;not null"`
	Enabled      bool   `json:"enabled" gorm:"not null;default:true"`
	LastStatus   string `json:"lastStatus" gorm:"size:30;not null;default:pending"`
	LastError    string `json:"lastError" gorm:"size:1000"`
	LastRunAt    *int64 `json:"lastRunAt"`
}

func (BackupPlan) TableName() string { return "backup_plan" }
