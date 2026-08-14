package main

import (
	"flag"
	"handfree-work/web-restic/internal/handler"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
)

var (
	port = flag.String("port", ":3000", "Port to listen on")
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// Connected with database

	// Create fiber app
	app := fiber.New(fiber.Config{})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// Create a /api/v1 endpoint
	handler.UserHandlersRegister(app)

	// Setup static files
	app.Get("/*", static.New("./static/public"))

	// Handle not founds
	// app.Use(handlers.NotFound)

	// Listen on port 3000
	log.Fatal(app.Listen(*port, fiber.ListenConfig{EnablePrefork: *prod})) // go run app.go -port=:3000
}
