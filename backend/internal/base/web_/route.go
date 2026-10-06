package web_

import "github.com/gofiber/fiber/v3"

// Route describes a POST business route and its optional permission.
// Additional route metadata can be added here without changing each handler's registration loop.
type Route struct {
	Path       string
	Permission string
	Handler    fiber.Handler
}

// RegisterRoutes registers routes and applies permission middleware when configured.
func RegisterRoutes(router fiber.Router, config Config, routes ...Route) {
	for _, route := range routes {
		handlers := []any{route.Handler}
		if route.Permission != "" {
			handlers = []any{Require(config, route.Permission), route.Handler}
		}
		router.Post(route.Path, handlers[0], handlers[1:]...)
	}
}
