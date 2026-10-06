package handler

import (
	"handfree-work/octo-backup/internal/base/log_"
	"handfree-work/octo-backup/internal/base/web_"
	"handfree-work/octo-backup/internal/models"
	"handfree-work/octo-backup/internal/svc"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// auditMiddleware 在认证请求完成后持久化用户操作，数据库写入失败只记录告警，不影响业务响应。
func auditMiddleware(svcCtx *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		startedAt := time.Now()
		err := c.Next()
		claims, authenticated := web_.ClaimsFromContext(c)
		if svcCtx.Db != nil && authenticated && len(c.Path()) >= 5 && c.Path()[:5] == "/api/" {
			entry := models.AuditLog{
				UserId: claims.UserId, Username: claims.Username, Method: c.Method(), Path: c.Path(),
				Status: c.Response().StatusCode(), Ip: c.IP(), Duration: time.Since(startedAt).Milliseconds(),
			}
			if createErr := svcCtx.Db.Create(&entry).Error; createErr != nil {
				log_.Logger.Warn("写入审计日志失败", zap.Error(createErr), zap.Int64("userId", claims.UserId), zap.String("path", c.Path()))
			}
		}
		return err
	}
}
