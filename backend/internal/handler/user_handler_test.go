package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"handfree-work/web-restic/internal/base/db_"
	"handfree-work/web-restic/internal/handler"
	"handfree-work/web-restic/internal/models"
	"handfree-work/web-restic/internal/svc"

	"github.com/gofiber/fiber/v3"
)

func TestUserCRUD(t *testing.T) {
	db, err := db_.OpenSQLite(filepath.Join(t.TempDir(), "db.sqlite"))
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("DB() error = %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	if err := db_.Migrate(db, &models.User{}); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}
	app := handler.NewApp(&svc.ServiceContext{Db: db})

	created := doJSONRequest(t, app, http.MethodPost, "/api/users", map[string]string{
		"username": "alice",
		"password": "secret",
		"nickName": "Alice",
	})
	if created.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d", created.Code, http.StatusCreated)
	}
	if created.Data["username"] != "alice" {
		t.Fatalf("created username = %#v, want alice", created.Data["username"])
	}
	if _, exists := created.Data["password"]; exists {
		t.Fatal("created response must not contain password")
	}
	userID := int64(created.Data["id"].(float64))

	list := doJSONRequest(t, app, http.MethodGet, "/api/users", nil)
	if list.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", list.Code, http.StatusOK)
	}
	users, ok := list.Data["items"].([]any)
	if !ok || len(users) != 1 {
		t.Fatalf("list items = %#v, want one user", list.Data["items"])
	}

	userPath := fmt.Sprintf("/api/users/%d", userID)
	found := doJSONRequest(t, app, http.MethodGet, userPath, nil)
	if found.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d", found.Code, http.StatusOK)
	}

	updated := doJSONRequest(t, app, http.MethodPut, userPath, map[string]string{
		"nickName": "Alice Updated",
	})
	if updated.Code != http.StatusOK {
		t.Fatalf("update status = %d, want %d", updated.Code, http.StatusOK)
	}
	if updated.Data["nickName"] != "Alice Updated" {
		t.Fatalf("updated nickname = %#v, want Alice Updated", updated.Data["nickName"])
	}

	deleted := doJSONRequest(t, app, http.MethodDelete, userPath, nil)
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d", deleted.Code, http.StatusNoContent)
	}

	notFound := doJSONRequest(t, app, http.MethodGet, userPath, nil)
	if notFound.Code != http.StatusNotFound {
		t.Fatalf("get deleted user status = %d, want %d", notFound.Code, http.StatusNotFound)
	}

	invalid := doJSONRequest(t, app, http.MethodPost, "/api/users", map[string]string{
		"password": "secret",
	})
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("invalid create status = %d, want %d", invalid.Code, http.StatusBadRequest)
	}

}

type response struct {
	Code int
	Data map[string]any
}

func doJSONRequest(t *testing.T, app *fiber.App, method, path string, body any) response {
	t.Helper()
	var payload io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request: %v", err)
		}
		payload = bytes.NewReader(encoded)
	}
	req := httptest.NewRequest(method, path, payload)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("request %s %s: %v", method, path, err)
	}
	t.Cleanup(func() { _ = res.Body.Close() })

	decoded := map[string]any{}
	if res.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
			t.Fatalf("decode response: %v", err)
		}
	}
	data, _ := decoded["data"].(map[string]any)
	return response{Code: res.StatusCode, Data: data}
}
