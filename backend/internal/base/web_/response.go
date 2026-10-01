package web_

import "github.com/gofiber/fiber/v3"

// Success returns the standard business response envelope.
func Success(c fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": 0, "message": "成功", "data": data})
}

// Error returns a business error while keeping the HTTP status at 200.
func Error(c fiber.Ctx, code int, message string) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": code, "message": message, "data": fiber.Map{}})
}

func BusinessError(c fiber.Ctx, err error) error {
	return Error(c, 1, err.Error())
}
