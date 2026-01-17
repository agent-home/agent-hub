package config

import (
	"os"
	"strconv"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Storage  StorageConfig
	Auth     AuthConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Address string
	Mode    string // debug, release, test
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// StorageConfig 对象存储配置
type StorageConfig struct {
	Type      string // local, s3, cos
	LocalPath string
	Bucket    string
	Region    string
	Endpoint  string
}

// AuthConfig 认证配置
type AuthConfig struct {
	JWTSecret     string
	TokenExpiry   int // 小时
	RefreshExpiry int // 天
}

// Load 加载配置
func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Address: getEnv("SERVER_ADDRESS", ":8080"),
			Mode:    getEnv("SERVER_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "agenthub"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "agenthub"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Storage: StorageConfig{
			Type:      getEnv("STORAGE_TYPE", "local"),
			LocalPath: getEnv("STORAGE_LOCAL_PATH", "./data/agents"),
			Bucket:    getEnv("STORAGE_BUCKET", ""),
			Region:    getEnv("STORAGE_REGION", ""),
			Endpoint:  getEnv("STORAGE_ENDPOINT", ""),
		},
		Auth: AuthConfig{
			JWTSecret:     getEnv("JWT_SECRET", "change-me-in-production"),
			TokenExpiry:   getEnvInt("TOKEN_EXPIRY_HOURS", 24),
			RefreshExpiry: getEnvInt("REFRESH_EXPIRY_DAYS", 7),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
