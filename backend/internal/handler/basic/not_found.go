package basic

import (
	"github.com/gofiber/fiber/v3"
)

func BaseHandlersRegister(app *fiber.App) {
	app.Use(NotFound)
}

// NotFound returns custom 404 page
func NotFound(c fiber.Ctx) error {
	return c.Status(404).SendFile("./static/private/404.html")
}
