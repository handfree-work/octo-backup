package logic

import (
	"context"
	provider "handfree-work/octo-backup/internal/modules/backup/restic"
	"testing"
)

func TestResticClientProvidesBackupAndRestoreOperations(t *testing.T) {
	var core interface {
		Backup(context.Context, provider.BackupRequest) (*provider.BackupResult, error)
		Restore(context.Context, provider.RestoreRequest) (string, error)
	} = &provider.ResticClient{}
	if core == nil {
		t.Fatal("ResticClient is nil")
	}
}
