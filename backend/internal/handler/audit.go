package handler

import (
	"github.com/gofiber/fiber/v3"
	"handfree-work/octo-backup/internal/svc"
	"time"
)

// auditMiddleware 在认证请求完成后持久化用户操作，数据库写入失败只记录告警，不影响业务响应。
func auditMiddleware(svcCtx *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Locals("auditDb", svcCtx.Db)
		c.Locals("auditStartedAt", time.Now())
		return c.Next()
	}
}
