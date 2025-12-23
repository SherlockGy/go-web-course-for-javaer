// 01-basic-config: Viper 基础配置读取
//
// 📌 Viper 功能:
//   - 支持多种配置格式 (YAML/JSON/TOML)
//   - 支持环境变量
//   - 支持配置热加载
//   - 支持远程配置
//
// 📌 最佳实践:
//   - 配置文件放在项目根目录或 configs/
//   - 使用环境变量覆盖配置
//   - 为配置项设置默认值
package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func main() {
	// 设置配置文件名（不带扩展名）
	viper.SetConfigName("config")

	// 设置配置文件类型
	viper.SetConfigType("yaml")

	// 添加配置文件搜索路径
	viper.AddConfigPath(".")         // 当前目录
	viper.AddConfigPath("./configs") // configs 目录

	// 设置默认值
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("database.max_connections", 10)

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("配置文件未找到，使用默认值")
		} else {
			log.Fatalf("读取配置文件失败: %v", err)
		}
	} else {
		log.Printf("使用配置文件: %s", viper.ConfigFileUsed())
	}

	// 读取配置值
	fmt.Println("=== 服务器配置 ===")
	fmt.Printf("Host: %s\n", viper.GetString("server.host"))
	fmt.Printf("Port: %d\n", viper.GetInt("server.port"))

	fmt.Println("\n=== 数据库配置 ===")
	fmt.Printf("DSN: %s\n", viper.GetString("database.dsn"))
	fmt.Printf("Max Connections: %d\n", viper.GetInt("database.max_connections"))

	fmt.Println("\n=== JWT 配置 ===")
	fmt.Printf("Secret: %s\n", viper.GetString("jwt.secret"))
	fmt.Printf("Expiration: %d 秒\n", viper.GetInt("jwt.expiration"))

	// 检查配置是否存在
	if viper.IsSet("logging.level") {
		fmt.Printf("\n日志级别: %s\n", viper.GetString("logging.level"))
	}

	// 获取所有配置
	fmt.Println("\n=== 所有配置 ===")
	for k, v := range viper.AllSettings() {
		fmt.Printf("%s: %v\n", k, v)
	}
}
