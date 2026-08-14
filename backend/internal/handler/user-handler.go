package handler

import (
	"handfree-work/web-restic/models"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/utils/v2"
)

func UserHandlersRegister(app *fiber.App) {
	v1 := app.Group("/api/user")
	// Bind handlers
	v1.Post("/list", UserList)
	v1.Post("/add", UserCreate)
}

// UserList returns a list of users
func UserList(c fiber.Ctx) error {
	users := []models.User{}

	return c.JSON(fiber.Map{
		"success": true,
		"users":   users,
	})
}

// UserCreate registers a user
func UserCreate(c fiber.Ctx) error {
	user := &models.User{
		// Note: when writing to external database,
		// we can simply use - Name: c.FormValue("user")
		Name: utils.CopyString(c.FormValue("user")),
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user":    user,
	})
}
