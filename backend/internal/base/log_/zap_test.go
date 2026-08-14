package log_

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestConsoleEncoderUsesCompactColoredPrefix(t *testing.T) {
	entryTime := time.Date(2026, time.August, 15, 0, 42, 4, 344_000_000, time.Local)
	tests := []struct {
		level zapcore.Level
		color string
		name  string
	}{
		{level: zap.InfoLevel, color: ansiBlue, name: "INFO"},
		{level: zap.WarnLevel, color: ansiYellow, name: "WARN"},
		{level: zap.ErrorLevel, color: ansiRed, name: "ERROR"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := newConsoleEncoder().EncodeEntry(zapcore.Entry{
				Level:   tt.level,
				Time:    entryTime,
				Message: "HTTP 请求",
			}, []zapcore.Field{zap.String("path", "/swagger/doc.json")})
			if err != nil {
				t.Fatalf("EncodeEntry() error = %v", err)
			}
			defer encoded.Free()
			output := encoded.String()
			wantPrefix := ansiCyan + "[2026-08-15 00:42:04.344]" + ansiReset + " " + tt.color + "[" + tt.name
			if !strings.Contains(output, wantPrefix) {
				t.Fatalf("console output = %q, want prefix containing %q", output, wantPrefix)
			}
			if strings.Contains(output, "\t") {
				t.Fatalf("console output must not contain tabs: %q", output)
			}
		})
	}
}

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
	if bytes.Contains(content, []byte("\x1b[")) {
		t.Fatalf("file log must not contain ANSI colors: %q", content)
	}
}
