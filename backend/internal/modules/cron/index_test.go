package cron

import "testing"

func TestServiceRegisterReplacesTaskByTypeAndId(t *testing.T) {
	service := New()
	defer service.Stop()
	if err := service.Register("backup_plan", 7, "0 0 * * *", func() {}); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if err := service.Register("backup_plan", 7, "5 0 * * *", func() {}); err != nil {
		t.Fatalf("replace failed: %v", err)
	}
	if len(service.entries) != 1 {
		t.Fatalf("expected one task after replacement, got %d", len(service.entries))
	}
}
