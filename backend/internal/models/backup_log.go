package models

import "handfree-work/octo-backup/internal/base/db_"

type BackupLog struct {
	db_.BaseModel
	PlanId     int64  `json:"planId" gorm:"not null;index"`
	Status     string `json:"status" gorm:"size:30;not null;index"`
	Progress   int    `json:"progress" gorm:"not null;default:0"`
	Stage      string `json:"stage" gorm:"size:100"`
	Result     string `json:"result" gorm:"type:text"`
	Error      string `json:"error" gorm:"type:text"`
	StartedAt  int64  `json:"startedAt"`
	FinishedAt *int64 `json:"finishedAt"`
}

func (BackupLog) TableName() string { return "backup_log" }
