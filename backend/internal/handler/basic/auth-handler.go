package basic

import (
	"handfree-work/web-restic/internal/auth"
	"handfree-work/web-restic/internal/base/web_"
	logic "handfree-work/web-restic/internal/service"
	"handfree-work/web-restic/internal/svc"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuth(app *fiber.App, svcCtx *svc.ServiceContext) {
	routes := app.Group("/api/auth")
	routes.Post("/register", auth.Require(svcCtx.Auth, auth.Guest), registerUser(svcCtx))
	routes.Post("/login", auth.Require(svcCtx.Auth, auth.Guest), loginUser(svcCtx))
}
func writeError(c fiber.Ctx, message string) error { return web_.Error(c, 1, message) }

// registerUser godoc
// @Summary 注册用户
// @Description 注册首个用户会自动授予 admin 角色，后续注册用户为 read 角色。
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body logic.CreateUserInput true "注册信息"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/auth/register [post]
func registerUser(svcCtx *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		var input logic.CreateUserInput
		if err := c.Bind().Body(&input); err != nil {
			return writeError(c, "请求体格式错误")
		}
		user, err := logic.NewUserService(c.Context(), svcCtx).Register(&input)
		if err != nil {
			return web_.BusinessError(c, err)
		}
		return web_.Success(c, user)
	}
}

// loginUser godoc
// @Summary 用户登录
// @Description 登录成功后返回有效期默认 7 天的 JWT Bearer token。
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body logic.LoginInput true "登录信息"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/login [post]
func loginUser(svcCtx *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		var input logic.LoginInput
		if err := c.Bind().Body(&input); err != nil {
			return writeError(c, "请求体格式错误")
		}
		user, err := logic.NewUserService(c.Context(), svcCtx).Login(&input)
		if err != nil {
			return web_.BusinessError(c, err)
		}
		if user.Id == nil {
			return writeError(c, "用户数据错误")
		}
		token, expiresAt, err := svcCtx.Auth.Issue(*user.Id, user.Username, user.Role)
		if err != nil {
			return writeError(c, "签发登录凭证失败")
		}
		return web_.Success(c, fiber.Map{
			"token":     token,
			"expiresAt": expiresAt.Unix(),
			"user": fiber.Map{
				"id":       user.Id,
				"username": user.Username,
				"role":     user.Role,
			},
		})
	}
}
