package models

// SysSetting 保存不对外暴露的系统级键值配置。
type SysSetting struct {
	Id      int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Key     string `json:"key" gorm:"column:key;size:100;not null;uniqueIndex"`
	Setting string `json:"setting" gorm:"column:setting;type:text;not null"`
}

// TableName 固定使用首版约定的单数表名。
func (SysSetting) TableName() string {
	return "sys_setting"
}
