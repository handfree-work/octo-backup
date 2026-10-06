package models

// AuditLog 记录已认证用户的 HTTP 操作及结果。
type AuditLog struct {
	Id        *int64 `json:"id" gorm:"primaryKey;column:id"`
	UserId    int64  `json:"userId" gorm:"column:user_id;index"`
	Username  string `json:"username" gorm:"column:username;size:100"`
	Operation string `json:"operation" gorm:"column:operation;size:100"`
	Method    string `json:"method" gorm:"column:method;size:10"`
	Path      string `json:"path" gorm:"column:path;size:255"`
	Status    int    `json:"status" gorm:"column:status"`
	Ip        string `json:"ip" gorm:"column:ip;size:64"`
	Duration  int64  `json:"duration" gorm:"column:duration;comment:耗时毫秒"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime:milli;column:created_at;index"`
}

// TableName 返回审计日志表名。
func (AuditLog) TableName() string { return "audit_log" }
