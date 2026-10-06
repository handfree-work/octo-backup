package handler

import (
	"handfree-work/octo-backup/internal/base/web_"
	"handfree-work/octo-backup/internal/handler/basic"
	sysHandler "handfree-work/octo-backup/internal/handler/sys"
	userHandler "handfree-work/octo-backup/internal/handler/user"
	"handfree-work/octo-backup/internal/svc"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

// NewApp 创建用于测试及嵌入场景的完整 HTTP 应用。
func NewApp(svcCtx *svc.ServiceContext) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: web_.ErrorHandler})
	app.Use(recover.New())
	Register(app, svcCtx)
	return app
}

// Register 按 basic、user、sys 模块统一注册 handler。
func Register(app *fiber.App, svcCtx *svc.ServiceContext) {
	app.Use(auditMiddleware(svcCtx))
	basic.Register(app)
	basic.RegisterAuth(app, svcCtx)
	userHandler.Register(app, svcCtx)
	sysHandler.Register(app, svcCtx)
}
