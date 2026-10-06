package models

import "handfree-work/octo-backup/internal/base/db_"

// RepoDataInfo 保存仓库详情的快照、统计和最近一次备份摘要缓存。
type RepoDataInfo struct {
	db_.BaseModel
	RepositoryId *int64 `json:"repositoryId" gorm:"column:repository_id;uniqueIndex;not null"`
	DataYAML     string `json:"-" gorm:"column:data_yaml;type:text;not null"`
}

// TableName 返回仓库详情缓存表名。
func (RepoDataInfo) TableName() string { return "repo_data_info" }
