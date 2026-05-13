package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config 全局配置结构
type Config struct {
	Server ServerConfig `mapstructure:"server"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	// 1. 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using system env vars")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv() // 允许环境变量覆盖 YAML 配置

	// 2. 读取默认配置文件（如果存在）
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config

	// 3. 手动绑定环境变量到结构体字段
	cfg.Server.Port = getEnvOrDefault("SERVER_PORT", "9000")

	return &cfg, nil
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
