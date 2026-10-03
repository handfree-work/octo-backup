package admin

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"handfree-work/octo-backup/internal/models"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("用户不存在")
	ErrInvalidInput = errors.New("用户名和密码不能为空")
)

// ResetPassword 使用 bcrypt 更新本地管理界面中的用户密码。
func ResetPassword(db *gorm.DB, username, password string) error {
	username = strings.TrimSpace(username)
	if username == "" || strings.TrimSpace(password) == "" {
		return ErrInvalidInput
	}
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: %s", ErrUserNotFound, username)
		}
		return fmt.Errorf("查询用户: %w", err)
	}
	digest := sha256.Sum256([]byte(password))
	if err := user.EncryptPassword(hex.EncodeToString(digest[:])); err != nil {
		return fmt.Errorf("加密密码: %w", err)
	}
	if err := db.Model(&models.User{}).Where("id = ?", user.Id).Update("password", user.Password).Error; err != nil {
		return fmt.Errorf("保存密码: %w", err)
	}
	return nil
}

// Run 启动 Bubble Tea 终端界面。
func Run(db *gorm.DB, input io.Reader, output io.Writer) error {
	return runTUI(db, input, output)
}

