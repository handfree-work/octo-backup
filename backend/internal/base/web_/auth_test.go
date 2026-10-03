package web_

import (
	"testing"
	"time"
)

func TestNewConfigDefaultsTokenTTLToSevenDays(t *testing.T) {
	config, err := NewConfig("test-secret", "")
	if err != nil {
		t.Fatalf("NewConfig() error = %v", err)
	}
	if config.TokenTTL != 7*24*time.Hour {
		t.Fatalf("TokenTTL = %s, want 168h", config.TokenTTL)
	}

	token, _, err := config.Issue(1, "alice", RoleRead)
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}
	claims, err := config.Parse(token)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if claims.UserID != 1 || claims.Role != RoleRead {
		t.Fatalf("claims = %#v, want user 1 with read role", claims)
	}
}

