package db_

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// OpenSQLite 打开 SQLite 数据库，并在需要时创建数据库目录。
func OpenSQLite(path string) (*gorm.DB, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("SQLite 数据库路径不能为空")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("创建数据库目录: %w", err)
	}
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, fmt.Errorf("打开 SQLite 数据库: %w", err)
	}
	return db, nil
}

// Migrate 迁移传入的数据库模型。
func Migrate(db *gorm.DB, models ...any) error {
	if db == nil {
		return fmt.Errorf("数据库连接不能为空")
	}
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("迁移数据库模型: %w", err)
	}
	return nil
}
