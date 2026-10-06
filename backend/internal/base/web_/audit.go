package web_

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"handfree-work/octo-backup/internal/base/log_"
	"handfree-work/octo-backup/internal/models"
	"time"
)

// AuditLog 写入当前业务操作的审计记录，身份和请求信息均从请求上下文读取。
func AuditLog(c fiber.Ctx, operation string) {
	db, ok := c.Locals("auditDb").(*gorm.DB)
	if !ok || db == nil || operation == "" {
		return
	}
	entry := models.AuditLog{
		Operation: operation,
		Method:    c.Method(),
		Path:      c.Path(),
		Status:    c.Response().StatusCode(),
		Ip:        c.IP(),
		Duration:  time.Since(auditStartedAt(c)).Milliseconds(),
	}
	if claims, authenticated := ClaimsFromContext(c); authenticated {
		entry.UserId, entry.Username = claims.UserId, claims.Username
	} else {
		if username, ok := c.Locals("auditUsername").(string); ok {
			entry.Username = username
		}
	}
	if err := db.Create(&entry).Error; err != nil {
		log_.Logger.Warn("写入审计日志失败", zap.Error(err), zap.String("operation", operation), zap.String("path", c.Path()))
	}
}

func auditStartedAt(c fiber.Ctx) time.Time {
	if startedAt, ok := c.Locals("auditStartedAt").(time.Time); ok {
		return startedAt
	}
	return time.Now()
}
