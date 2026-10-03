package sys

import (
	"handfree-work/octo-backup/internal/base/web_"
	logic "handfree-work/octo-backup/internal/service"
	"handfree-work/octo-backup/internal/svc"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func RegisterPlugin(app *fiber.App, s *svc.ServiceContext) {
	r := app.Group("/api/plugin")
	web_.RegisterRoutes(r, s.Auth,
		web_.Route{Path: "/metadata", Permission: web_.Read, Handler: pluginMetadata(s)},
		web_.Route{Path: "/page", Permission: web_.Read, Handler: pluginPage(s)},
		web_.Route{Path: "/create", Permission: web_.Write, Handler: pluginCreate(s)},
		web_.Route{Path: "/update", Permission: web_.Write, Handler: pluginUpdate(s)},
		web_.Route{Path: "/action", Permission: web_.Write, Handler: pluginAction(s)},
		web_.Route{Path: "/delete", Permission: web_.Admin, Handler: pluginDelete(s)},
	)
}
func pluginUpdate(s *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Query("id"), 10, 64)
		if err != nil {
			return web_.BusinessError(c, err)
		}
		var in logic.PluginInput
		if err = c.Bind().Body(&in); err != nil {
			return web_.BusinessError(c, err)
		}
		v, err := logic.NewPluginService(c.Context(), s).Update(id, &in)
		if err != nil {
			return web_.BusinessError(c, err)
		}
		return web_.Success(c, v)
	}
}
func pluginMetadata(s *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		var q struct {
			Type string `json:"type"`
			Name string `json:"name"`
		}
		if err := c.Bind().Body(&q); err != nil {
			return web_.BusinessError(c, err)
		}
		return web_.Success(c, logic.NewPluginService(c.Context(), s).Metadata(q.Type, q.Name))
	}
}
func pluginPage(s *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		var q logic.PluginPageQuery
		if err := c.Bind().Body(&q); err != nil {
			return web_.BusinessError(c, err)
		}
		v, err := logic.NewPluginService(c.Context(), s).Page(&q)
		if err != nil {
			return web_.BusinessError(c, err)
		}
		return web_.Success(c, v)
	}
}
func pluginCreate(s *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		var in logic.PluginInput
		if err := c.Bind().Body(&in); err != nil {
			return web_.BusinessError(c, err)
		}
		v, err := logic.NewPluginService(c.Context(), s).Create(&in)
		if err != nil {
			return web_.BusinessError(c, err)
		}
		return web_.Success(c, v)
	}
}
func pluginAction(s *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		var in struct {
			ID     int64          `json:"id"`
			Action string         `json:"action"`
			Params map[string]any `json:"params"`
		}
		if err := c.Bind().Body(&in); err != nil {
			return web_.BusinessError(c, err)
		}
		v, err := logic.NewPluginService(c.Context(), s).Action(in.ID, in.Action, in.Params)
		if err != nil {
			return web_.BusinessError(c, err)
		}
		return web_.Success(c, v)
	}
}
func pluginDelete(s *svc.ServiceContext) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Query("id"), 10, 64)
		if err != nil {
			return web_.BusinessError(c, err)
		}
		if err := logic.NewPluginService(c.Context(), s).Delete(id); err != nil {
			return web_.BusinessError(c, err)
		}
		return web_.Success(c, fiber.Map{})
	}
}

