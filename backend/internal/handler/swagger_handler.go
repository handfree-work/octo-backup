package handler

import (
	"errors"
	"io/fs"
	"mime"
	"path"
	"strings"

	docs "handfree-work/web-restic/docs"

	"github.com/gofiber/fiber/v3"
	swaggerFiles "github.com/swaggo/files/v2"
)

const swaggerIndex = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Web Restic API 文档</title>
  <link rel="stylesheet" href="./swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="./swagger-ui-bundle.js"></script>
  <script src="./swagger-ui-standalone-preset.js"></script>
  <script>
    window.ui = SwaggerUIBundle({
      url: "./doc.json",
      dom_id: "#swagger-ui",
      deepLinking: true,
      presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
      layout: "StandaloneLayout"
    });
  </script>
</body>
</html>`

// SwaggerHandlersRegister 注册 Swagger UI 与 OpenAPI JSON 文档。
func SwaggerHandlersRegister(app *fiber.App) {
	app.Get("/swagger", func(c fiber.Ctx) error {
		return c.Redirect().Status(fiber.StatusMovedPermanently).To("/swagger/index.html")
	})
	app.Get("/swagger/", func(c fiber.Ctx) error {
		return c.Redirect().Status(fiber.StatusMovedPermanently).To("/swagger/index.html")
	})
	app.Get("/swagger/doc.json", func(c fiber.Ctx) error {
		return c.Type("json").SendString(docs.SwaggerInfo.ReadDoc())
	})
	app.Get("/swagger/index.html", func(c fiber.Ctx) error {
		return c.Type("html").SendString(swaggerIndex)
	})
	app.Get("/swagger/*", swaggerAsset)
}

func swaggerAsset(c fiber.Ctx) error {
	fileName := path.Clean(strings.TrimPrefix(c.Params("*"), "/"))
	if fileName == "." || strings.HasPrefix(fileName, "../") {
		return c.SendStatus(fiber.StatusNotFound)
	}
	content, err := fs.ReadFile(swaggerFiles.FS, fileName)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return err
	}
	if contentType := mime.TypeByExtension(path.Ext(fileName)); contentType != "" {
		c.Set(fiber.HeaderContentType, contentType)
	}
	return c.Send(content)
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
