package sys

import (
	"handfree-work/octo-backup/internal/base/web_"
	logic "handfree-work/octo-backup/internal/modules/audit"
	"handfree-work/octo-backup/internal/svc"

	"github.com/gofiber/fiber/v3"
)

// RegisterAudit 注册审计日志查询接口。
func RegisterAudit(app *fiber.App, s *svc.ServiceContext) {
	r := app.Group("/api/audit")
	web_.RegisterRoutes(r, s.Auth, web_.Route{Path: "/page", Permission: web_.Read, Handler: web_.HandleJSON(func(c fiber.Ctx, q *logic.AuditLogPageQuery) (any, error) {
		return logic.NewService(c.Context(), s).Page(q)
	})})
}
