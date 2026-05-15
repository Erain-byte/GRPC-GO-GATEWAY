package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config 全局配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Services ServicesConfig `mapstructure:"services"`
	Consul   ConsulConfig   `mapstructure:"consul"`
}

// ... (其他结构体保持不变)

// ConsulConfig Consul 集群配置
type ConsulConfig struct {
	Addr       string `mapstructure:"addr"`
	Datacenter string `mapstructure:"datacenter"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// ServicesConfig 后端微服务发现配置
type ServicesConfig struct {
	AdminServiceName string `mapstructure:"admin_service_name"` // Consul 中的服务名
	// UserServicename  string `mapstructure:"user_service_name"`
	// AiServiceName    string `mapstructure:"ai_service_name"`
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
	cfg.Services.AdminServiceName = getEnvOrDefault("ADMIN_SERVICE_NAME", "admin-service")

	// Consul 配置
	cfg.Consul.Addr = getEnvOrDefault("CONSUL_ADDR", "127.0.0.1:8500")
	cfg.Consul.Datacenter = getEnvOrDefault("CONSUL_DATACENTER", "dc1")

	return &cfg, nil
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
