package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config 配置结构体
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	GRPC     GRPCConfig
	Log      LogConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Mode string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret string
	Expire string
}

// GRPCConfig gRPC 服务配置
type GRPCConfig struct {
	AIService   string
	UserService string
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string
	Format string
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: 未找到 .env 文件，使用环境变量")
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8081"),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			Username: getEnv("DB_USERNAME", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "admin_db"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "default-secret"),
			Expire: getEnv("JWT_EXPIRE", "24h"),
		},
		GRPC: GRPCConfig{
			AIService:   getEnv("GRPC_AI_SERVICE", "localhost:50052"),
			UserService: getEnv("GRPC_USER_SERVICE", "localhost:50051"),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
