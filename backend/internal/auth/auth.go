package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

const (
	Guest = "guest"
	Login = "login"
	Admin = "admin"
	Write = "write"
	Read  = "read"

	RoleAdmin = "admin"
	RoleWrite = "write"
	RoleRead  = "read"
)

var ErrInvalidToken = errors.New("无效的登录凭证")

type Config struct {
	Secret   string
	TokenTTL time.Duration
}

type Claims struct {
	UserID   int64  `json:"uid"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewConfig(secret, ttl string) (Config, error) {
	if strings.TrimSpace(secret) == "" {
		return Config{}, fmt.Errorf("JWT 密钥不能为空")
	}
	duration := 7 * 24 * time.Hour
	if strings.TrimSpace(ttl) != "" {
		parsed, err := time.ParseDuration(ttl)
		if err != nil || parsed <= 0 {
			return Config{}, fmt.Errorf("JWT 有效期无效")
		}
		duration = parsed
	}
	return Config{Secret: secret, TokenTTL: duration}, nil
}

func (c Config) Issue(userID int64, username, role string) (string, time.Time, error) {
	if strings.TrimSpace(c.Secret) == "" {
		return "", time.Time{}, fmt.Errorf("JWT 密钥不能为空")
	}
	ttl := c.TokenTTL
	if ttl <= 0 {
		ttl = 7 * 24 * time.Hour
	}
	expiresAt := time.Now().Add(ttl)
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(c.Secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("签发 JWT: %w", err)
	}
	return signed, expiresAt, nil
}

func (c Config) Parse(tokenText string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenText, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(c.Secret), nil
	})
	if err != nil || !token.Valid || !ValidRole(claims.Role) {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func Require(config Config, permission string) fiber.Handler {
	return func(c fiber.Ctx) error {
		if permission == Guest {
			return c.Next()
		}
		claims, err := config.Parse(bearerToken(c.Get(fiber.HeaderAuthorization)))
		if err != nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": fiber.StatusUnauthorized, "message": "未登录或登录已过期", "data": fiber.Map{}})
		}
		if !HasPermission(claims.Role, permission) {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": fiber.StatusForbidden, "message": "没有访问权限", "data": fiber.Map{}})
		}
		c.Locals("authClaims", claims)
		return c.Next()
	}
}

func ValidRole(role string) bool {
	return role == RoleAdmin || role == RoleWrite || role == RoleRead
}

func HasPermission(role, permission string) bool {
	if permission == Login {
		return ValidRole(role)
	}
	switch permission {
	case Admin:
		return role == RoleAdmin
	case Write:
		return role == RoleAdmin || role == RoleWrite
	case Read:
		return role == RoleAdmin || role == RoleWrite || role == RoleRead
	default:
		return false
	}
}

func ClaimsFromContext(c fiber.Ctx) (*Claims, bool) {
	claims, ok := c.Locals("authClaims").(*Claims)
	return claims, ok
}

func bearerToken(header string) string {
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return parts[1]
}
