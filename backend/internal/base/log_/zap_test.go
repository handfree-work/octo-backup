package log_

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitZapWritesToLogFile(t *testing.T) {
	directory := t.TempDir()
	logger, err := InitZap(ZapConfig{Mode: "dev", Directory: directory, Level: "debug"})
	if err != nil {
		t.Fatalf("InitZap() error = %v", err)
	}
	logger.Info("file log test")
	t.Cleanup(ToDefer)
	if err := logger.Sync(); err != nil {
		t.Logf("Sync() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(directory, "app.log"))
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}
	if !strings.Contains(string(content), "file log test") {
		t.Fatalf("log file = %q, want message", content)
	}
}
