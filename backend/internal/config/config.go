package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name     string         `yaml:"name"`
	Mode     string         `yaml:"mode"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

// Load 先加载 dev.yaml 作为默认配置，再加载指定运行模式的配置覆盖同名字段。
func Load(configDir, mode string) (*Config, error) {
	if strings.TrimSpace(mode) == "" {
		mode = "dev"
	}

	var cfg Config
	if err := loadFile(filepath.Join(configDir, "dev.yaml"), &cfg); err != nil {
		return nil, err
	}
	if mode != "dev" {
		if err := loadFile(filepath.Join(configDir, mode+".yaml"), &cfg); err != nil {
			return nil, err
		}
	}
	if cfg.Mode == "" {
		cfg.Mode = mode
	}
	return &cfg, nil
}

func loadFile(path string, cfg *Config) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取配置文件 %s: %w", path, err)
	}
	if err := yaml.Unmarshal(content, cfg); err != nil {
		return fmt.Errorf("解析配置文件 %s: %w", path, err)
	}
	return nil
}
