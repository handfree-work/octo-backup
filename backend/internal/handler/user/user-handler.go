package user

import (
	"errors"
	"strconv"

	"handfree-work/web-restic/internal/auth"
	logic "handfree-work/web-restic/internal/service"
	"handfree-work/web-restic/internal/svc"

	"github.com/gofiber/fiber/v3"
)

// Register 注册认证与用户管理接口。
func Register(app *fiber.App, svcCtx *svc.ServiceContext) {
	authRoutes := app.Group("/api/auth")
	authRoutes.Post("/register", auth.Require(svcCtx.Auth, auth.Guest), registerUser(svcCtx))
	authRoutes.Post("/login", auth.Require(svcCtx.Auth, auth.Guest), loginUser(svcCtx))

	users := app.Group("/api/users")
	users.Post("/create", auth.Require(svcCtx.Auth, auth.Admin), createUser(svcCtx))
	users.Post("/list", auth.Require(svcCtx.Auth, auth.Read), listUsers(svcCtx))
	users.Post("/:id/detail", auth.Require(svcCtx.Auth, auth.Read), getUser(svcCtx))
	users.Post("/:id/update", auth.Require(svcCtx.Auth, auth.Write), updateUser(svcCtx))
	users.Post("/:id/delete", auth.Require(svcCtx.Auth, auth.Admin), deleteUser(svcCtx))
}

// createUser godoc
// @Summary 创建用户
// @Description 需要 admin 权限，可设置 admin、write 或 read 角色。
// @Tags 用户
// @Accept json
// @Produce json
// @Security bearerAuth
// @Param request body logic.CreateUserInput true "用户信息"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/users/create [post]
func createUser(svcCtx *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		var input logic.CreateUserInput
		if err := c.Bind().Body(&input); err != nil {
			return writeError(c, fiber.StatusBadRequest, "请求体格式错误")
		}
		user, err := logic.NewUserService(c.Context(), svcCtx).Create(&input)
		if err != nil {
			return writeUserError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": user})
	}
}

// listUsers godoc
// @Summary 获取用户列表
// @Description 需要 read、write 或 admin 权限。
// @Tags 用户
// @Produce json
// @Security bearerAuth
// @Success 200 {object} UserListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /api/users/list [post]
func listUsers(svcCtx *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		users, err := logic.NewUserService(c.Context(), svcCtx).List()
		if err != nil {
			return writeUserError(c, err)
		}
		return c.JSON(fiber.Map{"data": fiber.Map{"items": users}})
	}
}

// getUser godoc
// @Summary 获取用户详情
// @Description 需要 read、write 或 admin 权限。
// @Tags 用户
// @Produce json
// @Security bearerAuth
// @Param id path int true "用户 ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id}/detail [post]
func getUser(svcCtx *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseUserID(c)
		if err != nil {
			return writeError(c, fiber.StatusBadRequest, "用户 ID 无效")
		}
		user, err := logic.NewUserService(c.Context(), svcCtx).Get(id)
		if err != nil {
			return writeUserError(c, err)
		}
		return c.JSON(fiber.Map{"data": user})
	}
}

// updateUser godoc
// @Summary 更新用户
// @Description 需要 write 或 admin 权限；非管理员只能更新自己的资料且不能修改角色。
// @Tags 用户
// @Accept json
// @Produce json
// @Security bearerAuth
// @Param id path int true "用户 ID"
// @Param request body logic.UpdateUserInput true "更新信息"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id}/update [post]
func updateUser(svcCtx *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseUserID(c)
		if err != nil {
			return writeError(c, fiber.StatusBadRequest, "用户 ID 无效")
		}
		var input logic.UpdateUserInput
		if err := c.Bind().Body(&input); err != nil {
			return writeError(c, fiber.StatusBadRequest, "请求体格式错误")
		}
		claims, ok := auth.ClaimsFromContext(c)
		if !ok {
			return writeError(c, fiber.StatusUnauthorized, "未登录或登录已过期")
		}
		if claims.Role != auth.RoleAdmin && claims.UserID != id {
			return writeError(c, fiber.StatusForbidden, "没有访问权限")
		}
		if claims.Role != auth.RoleAdmin && input.Role != nil {
			return writeError(c, fiber.StatusForbidden, "只有管理员可以调整角色")
		}
		user, err := logic.NewUserService(c.Context(), svcCtx).Update(id, &input)
		if err != nil {
			return writeUserError(c, err)
		}
		return c.JSON(fiber.Map{"data": user})
	}
}

// deleteUser godoc
// @Summary 删除用户
// @Description 需要 admin 权限。
// @Tags 用户
// @Produce json
// @Security bearerAuth
// @Param id path int true "用户 ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id}/delete [post]
func deleteUser(svcCtx *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := parseUserID(c)
		if err != nil {
			return writeError(c, fiber.StatusBadRequest, "用户 ID 无效")
		}
		if err := logic.NewUserService(c.Context(), svcCtx).Delete(id); err != nil {
			return writeUserError(c, err)
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func parseUserID(c fiber.Ctx) (int64, error) {
	return strconv.ParseInt(c.Params("id"), 10, 64)
}

func writeUserError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, logic.ErrInvalidUser):
		return writeError(c, fiber.StatusBadRequest, err.Error())
	case errors.Is(err, logic.ErrUserNotFound):
		return writeError(c, fiber.StatusNotFound, err.Error())
	case errors.Is(err, logic.ErrUsernameExists):
		return writeError(c, fiber.StatusConflict, err.Error())
	case errors.Is(err, logic.ErrInvalidLogin):
		return writeError(c, fiber.StatusUnauthorized, err.Error())
	default:
		return writeError(c, fiber.StatusInternalServerError, "服务器内部错误")
	}
}

func writeError(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"error": message})
}

// ErrorResponse 表示接口错误响应。
type ErrorResponse struct {
	Error string `json:"error" example:"没有访问权限"`
}

// UserResponse 表示单个用户的成功响应。
type UserResponse struct {
	Data map[string]any `json:"data"`
}

// UserListResponse 表示用户列表的成功响应。
type UserListResponse struct {
	Data struct {
		Items []map[string]any `json:"items"`
	} `json:"data"`
}

// LoginResponse 表示登录成功响应。
type LoginResponse struct {
	Data struct {
		Token     string         `json:"token"`
		ExpiresAt int64          `json:"expiresAt"`
		User      map[string]any `json:"user"`
	} `json:"data"`
}
