// 03-env-override: 环境变量覆盖配置
//
// 📌 环境变量覆盖的价值:
//   - 12-Factor App 最佳实践
//   - 不同环境使用不同配置
//   - 敏感信息不入代码库
//
// 📌 优先级（从低到高）:
//   1. 默认值
//   2. 配置文件
//   3. 环境变量
//   4. 命令行参数
//
// 测试:
//   APP_SERVER_PORT=9090 go run main.go
//   APP_DATABASE_DSN=":memory:" APP_JWT_SECRET="env-secret" go run main.go
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 🔑 关键：启用环境变量
	viper.SetEnvPrefix("APP") // 环境变量前缀，如 APP_SERVER_PORT
	viper.AutomaticEnv()      // 自动绑定环境变量

	// 将 . 替换为 _，使 server.port 对应 APP_SERVER_PORT
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 设置默认值
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.dsn", "./data.db")
	viper.SetDefault("jwt.secret", "default-secret")

	// 读取配置文件（可选）
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		log.Println("配置文件未找到，使用默认值和环境变量")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	fmt.Println("=== 当前配置 ===")
	fmt.Printf("Server Host: %s\n", cfg.Server.Host)
	fmt.Printf("Server Port: %d\n", cfg.Server.Port)
	fmt.Printf("Database DSN: %s\n", cfg.Database.DSN)
	fmt.Printf("JWT Secret: %s\n", cfg.JWT.Secret)

	fmt.Println("\n=== 测试环境变量覆盖 ===")
	fmt.Println("运行以下命令测试环境变量覆盖:")
	fmt.Println("  APP_SERVER_PORT=9090 go run main.go")
	fmt.Println("  APP_DATABASE_DSN=':memory:' APP_JWT_SECRET='new-secret' go run main.go")
}
