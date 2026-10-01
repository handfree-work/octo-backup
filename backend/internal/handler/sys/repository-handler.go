package sys

import (
	"github.com/gofiber/fiber/v3"
	"handfree-work/web-restic/internal/auth"
	"handfree-work/web-restic/internal/base/web_"
	logic "handfree-work/web-restic/internal/service"
	"handfree-work/web-restic/internal/svc"
	"strconv"
)

func RegisterRepository(app *fiber.App, svcCtx *svc.ServiceContext) {
	r := app.Group("/api/repository")
	r.Post("/page", auth.Require(svcCtx.Auth, auth.Read), repositoryPage(svcCtx))
	r.Post("/create", auth.Require(svcCtx.Auth, auth.Write), repositoryCreate(svcCtx))
	r.Post("/update", auth.Require(svcCtx.Auth, auth.Write), repositoryUpdate(svcCtx))
	r.Post("/delete", auth.Require(svcCtx.Auth, auth.Admin), repositoryDelete(svcCtx))
}
func repositoryPage(s *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		var q logic.StorageRepositoryPageQuery
		if err := c.Bind().Body(&q); err != nil {
			return web_.BusinessError(c, err)
		}
		v, err := logic.NewStorageRepositoryService(c.Context(), s).FindPage(&q)
		if err != nil {
			return web_.BusinessError(c, err)
		}
		return web_.Success(c, v)
	}
}
func repositoryCreate(s *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		var in logic.StorageRepositoryInput
		if err := c.Bind().Body(&in); err != nil {
			return web_.BusinessError(c, err)
		}
		v, err := logic.NewStorageRepositoryService(c.Context(), s).Create(&in)
		if err != nil {
			return web_.BusinessError(c, err)
		}
		return web_.Success(c, v)
	}
}
func repositoryUpdate(s *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Query("id"), 10, 64)
		if err != nil {
			return web_.BusinessError(c, err)
		}
		var in logic.UpdateStorageRepositoryInput
		if err = c.Bind().Body(&in); err != nil {
			return web_.BusinessError(c, err)
		}
		v, err := logic.NewStorageRepositoryService(c.Context(), s).Update(id, &in)
		if err != nil {
			return web_.BusinessError(c, err)
		}
		return web_.Success(c, v)
	}
}
func repositoryDelete(s *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Query("id"), 10, 64)
		if err != nil {
			return web_.BusinessError(c, err)
		}
		if err = logic.NewStorageRepositoryService(c.Context(), s).Delete(id); err != nil {
			return web_.BusinessError(c, err)
		}
		return web_.Success(c, fiber.Map{})
	}
}
