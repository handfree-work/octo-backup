package handler

import (
	"errors"
	"strconv"

	logic "handfree-work/web-restic/internal/service"
	"handfree-work/web-restic/internal/svc"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func NewApp(svcCtx *svc.ServiceContext) *fiber.App {
	app := fiber.New()
	app.Use(recover.New())
	UserHandlersRegister(app, svcCtx)
	return app
}

func UserHandlersRegister(app *fiber.App, svcCtx *svc.ServiceContext) {
	users := app.Group("/api/users")
	users.Post("/", createUser(svcCtx))
	users.Get("/", listUsers(svcCtx))
	users.Get("/:id", getUser(svcCtx))
	users.Put("/:id", updateUser(svcCtx))
	users.Delete("/:id", deleteUser(svcCtx))
}

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

func listUsers(svcCtx *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		users, err := logic.NewUserService(c.Context(), svcCtx).List()
		if err != nil {
			return writeUserError(c, err)
		}
		return c.JSON(fiber.Map{"data": fiber.Map{"items": users}})
	}
}

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
		user, err := logic.NewUserService(c.Context(), svcCtx).Update(id, &input)
		if err != nil {
			return writeUserError(c, err)
		}
		return c.JSON(fiber.Map{"data": user})
	}
}

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
	default:
		return writeError(c, fiber.StatusInternalServerError, "服务器内部错误")
	}
}

func writeError(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"error": message})
}
