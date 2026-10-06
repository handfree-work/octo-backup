package db_

import (
	"handfree-work/octo-backup/internal/base/error_"
	"os"
	"path/filepath"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// OpenSQLite 打开 SQLite 数据库，并在需要时创建数据库目录。
func OpenSQLite(path string) (*gorm.DB, error) {
	if strings.TrimSpace(path) == "" {
		return nil, error_.NewTextError("SQLite 数据库路径不能为空")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, error_.NewWrapError("创建数据库目录", err)
	}
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{TranslateError: true, Logger: logger.Discard})
	if err != nil {
		return nil, error_.NewWrapError("打开 SQLite 数据库", err)
	}
	return db, nil
}

// Migrate 迁移传入的数据库模型。
func Migrate(db *gorm.DB, models ...any) error {
	if db == nil {
		return error_.NewTextError("数据库连接不能为空")
	}
	if err := db.AutoMigrate(models...); err != nil {
		return error_.NewWrapError("迁移数据库模型", err)
	}
	return nil
}
