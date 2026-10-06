package models

import "handfree-work/octo-backup/internal/base/db_"

type PluginInstance struct {
	db_.BaseModel
	Name        string `json:"name" gorm:"size:100;not null"`
	PluginType  string `json:"pluginType" gorm:"size:30;not null;index"`
	PluginName  string `json:"pluginName" gorm:"size:100;not null;index"`
	ConfigYAML  string `json:"-" gorm:"type:text;not null"`
	Vars        string `json:"-" gorm:"column:vars;type:text;not null;default:''"`
	Description string `json:"description" gorm:"size:500"`
}

func (PluginInstance) TableName() string { return "plugin_instance" }
