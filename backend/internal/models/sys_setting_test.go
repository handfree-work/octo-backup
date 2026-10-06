package models_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"handfree-work/octo-backup/internal/base/db_"
	"handfree-work/octo-backup/internal/models"
)

func TestSysSettingSchemaAndPersistence(t *testing.T) {
	database, err := db_.OpenSQLite(filepath.Join(t.TempDir(), "db.sqlite"))
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		t.Fatalf("DB() error = %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	if err := db_.Migrate(database, &models.SysSetting{}); err != nil {
		t.Fatalf("Migrate() error = %v", err)
	}
	if !database.Migrator().HasTable("sys_setting") {
		t.Fatal("sys_setting table was not created")
	}
	columnTypes, err := database.Migrator().ColumnTypes("sys_setting")
	if err != nil {
		t.Fatalf("ColumnTypes() error = %v", err)
	}
	columns := make([]string, 0, len(columnTypes))
	for _, column := range columnTypes {
		columns = append(columns, column.Name())
	}
	wantColumns := []string{"id", "key", "setting"}
	if !reflect.DeepEqual(columns, wantColumns) {
		t.Fatalf("sys_setting columns = %#v, want %#v", columns, wantColumns)
	}

	setting := &models.SysSetting{Key: "instance", Setting: `{"name":"OctoBackup"}`}
	if err := database.Create(setting).Error; err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if setting.Id == 0 {
		t.Fatal("created setting must have an Id")
	}
	var found models.SysSetting
	if err := database.Where("key = ?", "instance").First(&found).Error; err != nil {
		t.Fatalf("First() error = %v", err)
	}
	if found.Setting != setting.Setting {
		t.Fatalf("stored setting = %q, want %q", found.Setting, setting.Setting)
	}
}
