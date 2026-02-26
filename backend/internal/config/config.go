package config

import (
	"os"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig HTTP服务配置
type ServerConfig struct {
	Port string
	Mode string // debug/release/test
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver string // sqlite/mysql
	DSN    string // 数据库连接字符串
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string
	ExpireHour int
}

// Load 从环境变量加载配置，使用默认值作为回退
func Load() *Config {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8000"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Driver: getEnv("DB_DRIVER", "sqlite"),
			DSN:    getEnv("DB_DSN", "marsview.db"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "marsview-secret-key-2024"),
			ExpireHour: 72,
		},
	}
	return cfg
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
