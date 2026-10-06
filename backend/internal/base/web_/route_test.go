package web_

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestRegisterRoutesPermission(t *testing.T) {
	app := fiber.New()
	RegisterRoutes(app, Config{Secret: "test-secret"},
		Route{Path: "/public", Handler: func(c fiber.Ctx) error { return c.SendStatus(http.StatusNoContent) }},
		Route{Path: "/admin", Permission: Read, Handler: func(c fiber.Ctx) error { return c.SendStatus(http.StatusNoContent) }},
	)

	public, err := app.Test(httptest.NewRequest(http.MethodPost, "/public", nil))
	if err != nil {
		t.Fatalf("POST /public: %v", err)
	}
	_ = public.Body.Close()
	if public.StatusCode != http.StatusNoContent {
		t.Fatalf("POST /public status = %d, want %d", public.StatusCode, http.StatusNoContent)
	}
	protected, err := app.Test(httptest.NewRequest(http.MethodPost, "/admin", nil))
	if err != nil {
		t.Fatalf("POST /admin: %v", err)
	}
	_ = protected.Body.Close()
	if protected.StatusCode != http.StatusOK {
		t.Fatalf("POST /admin status = %d, want %d", protected.StatusCode, http.StatusOK)
	}
}
