// 01-basic-zap: Zap 基础使用
//
// 📌 Logger vs SugaredLogger:
//   - Logger: 高性能，类型安全，适合生产环境
//   - SugaredLogger: 更灵活，支持 printf 风格，开发友好
//
// 📌 日志级别（从低到高）:
//   Debug → Info → Warn → Error → DPanic → Panic → Fatal
package main

import (
	"go.uber.org/zap"
)

func main() {
	// ==================== 快速创建 ====================

	// 开发模式：输出格式化、彩色、带调用栈
	devLogger, _ := zap.NewDevelopment()
	defer devLogger.Sync()

	devLogger.Info("开发模式日志")
	devLogger.Debug("Debug 信息")

	// 生产模式：JSON 格式，高性能
	prodLogger, _ := zap.NewProduction()
	defer prodLogger.Sync()

	prodLogger.Info("生产模式日志",
		zap.String("service", "user-api"),
		zap.Int("port", 8080),
	)

	// ==================== Logger（高性能）====================

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// 结构化日志字段
	logger.Info("用户登录",
		zap.Int64("user_id", 12345),
		zap.String("username", "tom"),
		zap.String("ip", "192.168.1.1"),
	)

	logger.Warn("请求超时",
		zap.String("endpoint", "/api/users"),
		zap.Duration("timeout", 5000000000), // 5秒
	)

	logger.Error("数据库连接失败",
		zap.String("dsn", "localhost:3306"),
		zap.Error(nil), // 可以传入 error
	)

	// ==================== SugaredLogger（更灵活）====================

	sugar := logger.Sugar()

	// printf 风格
	sugar.Infof("用户 %s 登录成功", "tom")
	sugar.Warnf("请求 %s 超时，耗时 %dms", "/api/users", 5000)

	// 键值对风格
	sugar.Infow("订单创建",
		"order_id", "ORD123",
		"user_id", 12345,
		"amount", 99.9,
	)

	// 自由格式
	sugar.Info("简单日志消息")
}
