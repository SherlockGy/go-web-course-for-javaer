// 02-struct-mapping: 配置映射到结构体
//
// 📌 结构体映射的好处:
//   - 类型安全
//   - IDE 自动补全
//   - 易于测试
//
// 📌 mapstructure tag:
//   - 用于映射配置键名到字段
//   - 支持嵌套结构体
package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	DSN            string `mapstructure:"dsn"`
	MaxConnections int    `mapstructure:"max_connections"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration int    `mapstructure:"expiration"`
}

type LoggingConfig struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

// LoadConfig 加载配置
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")

	// 设置默认值
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.max_connections", 10)
	viper.SetDefault("jwt.expiration", 86400)
	viper.SetDefault("logging.level", "info")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置失败: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	return &config, nil
}

func main() {
	cfg, err := LoadConfig(".")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 使用配置（类型安全）
	fmt.Printf("服务器地址: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("数据库 DSN: %s\n", cfg.Database.DSN)
	fmt.Printf("最大连接数: %d\n", cfg.Database.MaxConnections)
	fmt.Printf("JWT 过期时间: %d 秒\n", cfg.JWT.Expiration)
	fmt.Printf("日志级别: %s\n", cfg.Logging.Level)

	// 在实际项目中，通常将 config 作为全局变量或依赖注入
	// var globalConfig *Config
	// globalConfig = cfg
}
