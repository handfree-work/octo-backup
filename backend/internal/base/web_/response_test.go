package web_

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestErrorUsesDefaultCode(t *testing.T) {
	app := fiber.New()
	app.Post("/error", func(c fiber.Ctx) error {
		return Error(c, "请求失败")
	})

	res, err := app.Test(httptest.NewRequest("POST", "/error", nil))
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	defer res.Body.Close()

	var got struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Code != 1 || got.Message != "请求失败" {
		t.Fatalf("response = %#v, want code 1 and message", got)
	}
}

func TestBusinessErrorUsesProvidedCode(t *testing.T) {
	app := fiber.New()
	app.Post("/error", func(c fiber.Ctx) error {
		return BusinessError(c, 400, errors.New("参数错误"))
	})

	res, err := app.Test(httptest.NewRequest("POST", "/error", nil))
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	defer res.Body.Close()

	var got struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Code != 400 || got.Message != "参数错误" {
		t.Fatalf("response = %#v, want code 400 and message", got)
	}
}

func TestErrorHandlerReturnsBusinessEnvelope(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	app.Post("/error", func(fiber.Ctx) error { return errors.New("业务执行失败") })

	res, err := app.Test(httptest.NewRequest("POST", "/error", nil))
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	defer res.Body.Close()
	var got struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if res.StatusCode != fiber.StatusOK || got.Code == 0 || got.Message != "业务执行失败" {
		t.Fatalf("response status=%d body=%#v", res.StatusCode, got)
	}
}

func TestHandleJSONBindsBodyAndWrapsSuccess(t *testing.T) {
	type request struct {
		Name string `json:"name"`
	}
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	app.Post("/users", HandleJSON(func(_ fiber.Ctx, input *request) (any, error) {
		return fiber.Map{"name": input.Name}, nil
	}))

	bodyRequest := httptest.NewRequest("POST", "/users", strings.NewReader(`{"name":"alice"}`))
	bodyRequest.Header.Set("Content-Type", "application/json")
	res, err := app.Test(bodyRequest)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	defer res.Body.Close()
	var got struct {
		Code int `json:"code"`
		Data struct {
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Code != 0 || got.Data.Name != "alice" {
		t.Fatalf("response = %#v", got)
	}
}
