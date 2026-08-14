package handler_test

import (
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"handfree-work/web-restic/internal/auth"
	"handfree-work/web-restic/internal/base/db_"
	"handfree-work/web-restic/internal/handler"
	"handfree-work/web-restic/internal/models"
	"handfree-work/web-restic/internal/svc"
)

func TestRegisterLoginAndDeclaredPermissions(t *testing.T) {
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
	app := handler.NewApp(&svc.ServiceContext{
		Db:   db,
		Auth: auth.Config{Secret: "test-secret", TokenTTL: time.Hour},
	})

	unauthorized := doJSONRequest(t, app, http.MethodGet, "/api/users", nil)
	if unauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("unauthorized list status = %d, want %d", unauthorized.Code, http.StatusUnauthorized)
	}

	registered := doJSONRequest(t, app, http.MethodPost, "/api/auth/register", map[string]string{
		"username": "admin",
		"password": "secret",
	})
	if registered.Code != http.StatusCreated {
		t.Fatalf("register status = %d, want %d", registered.Code, http.StatusCreated)
	}
	if registered.Data["role"] != "admin" {
		t.Fatalf("first registered role = %#v, want admin", registered.Data["role"])
	}

	badLogin := doJSONRequest(t, app, http.MethodPost, "/api/auth/login", map[string]string{
		"username": "admin",
		"password": "wrong",
	})
	if badLogin.Code != http.StatusUnauthorized {
		t.Fatalf("bad login status = %d, want %d", badLogin.Code, http.StatusUnauthorized)
	}

	login := doJSONRequest(t, app, http.MethodPost, "/api/auth/login", map[string]string{
		"username": "admin",
		"password": "secret",
	})
	if login.Code != http.StatusOK {
		t.Fatalf("login status = %d, want %d", login.Code, http.StatusOK)
	}
	token, ok := login.Data["token"].(string)
	if !ok || token == "" {
		t.Fatalf("login token = %#v, want non-empty string", login.Data["token"])
	}

	authorized := doJSONRequest(t, app, http.MethodGet, "/api/users", nil, token)
	if authorized.Code != http.StatusOK {
		t.Fatalf("authorized list status = %d, want %d", authorized.Code, http.StatusOK)
	}

	second := doJSONRequest(t, app, http.MethodPost, "/api/auth/register", map[string]string{
		"username": "reader",
		"password": "secret",
	})
	if second.Code != http.StatusCreated || second.Data["role"] != "read" {
		t.Fatalf("second registered response = %#v, want created read user", second)
	}
	readerLogin := doJSONRequest(t, app, http.MethodPost, "/api/auth/login", map[string]string{
		"username": "reader",
		"password": "secret",
	})
	readerToken, ok := readerLogin.Data["token"].(string)
	if !ok || readerToken == "" {
		t.Fatalf("reader token = %#v, want non-empty string", readerLogin.Data["token"])
	}
	forbidden := doJSONRequest(t, app, http.MethodPut, "/api/users/2", map[string]string{
		"nickName": "blocked",
	}, readerToken)
	if forbidden.Code != http.StatusForbidden {
		t.Fatalf("reader update status = %d, want %d", forbidden.Code, http.StatusForbidden)
	}
}
