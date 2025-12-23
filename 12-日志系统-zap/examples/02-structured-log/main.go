// 02-structured-log: 结构化日志
//
// 📌 结构化日志的价值:
//   - 机器可解析（JSON）
//   - 便于日志聚合和分析
//   - 支持复杂查询
//
// 📌 常用字段类型:
//   zap.String, zap.Int, zap.Int64, zap.Float64
//   zap.Bool, zap.Time, zap.Duration
//   zap.Error, zap.Any, zap.Reflect
package main

import (
	"errors"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 自定义配置
func createLogger() *zap.Logger {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Encoding:    "json", // 或 "console"
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := config.Build()
	return logger
}

func main() {
	logger := createLogger()
	defer logger.Sync()

	// ==================== 请求日志示例 ====================
	logger.Info("HTTP 请求",
		zap.String("method", "POST"),
		zap.String("path", "/api/users"),
		zap.Int("status", 200),
		zap.Duration("latency", 45*time.Millisecond),
		zap.String("client_ip", "192.168.1.100"),
		zap.String("user_agent", "Mozilla/5.0"),
		zap.String("request_id", "req-abc-123"),
	)

	// ==================== 业务日志示例 ====================
	logger.Info("订单创建",
		zap.String("order_id", "ORD-2024-001"),
		zap.Int64("user_id", 12345),
		zap.Float64("amount", 199.99),
		zap.Strings("items", []string{"iPhone", "Case"}),
		zap.Time("created_at", time.Now()),
	)

	// ==================== 错误日志示例 ====================
	err := errors.New("connection refused")
	logger.Error("数据库连接失败",
		zap.String("host", "localhost"),
		zap.Int("port", 3306),
		zap.Int("retry_count", 3),
		zap.Error(err),
	)

	// ==================== 使用子日志器 ====================
	userLogger := logger.With(
		zap.String("module", "user"),
		zap.String("service", "user-api"),
	)

	userLogger.Info("用户注册", zap.String("username", "tom"))
	userLogger.Info("用户登录", zap.String("username", "tom"))

	// ==================== 复杂数据 ====================
	logger.Info("复杂数据",
		zap.Any("config", map[string]any{
			"host":     "localhost",
			"port":     8080,
			"features": []string{"auth", "rate-limit"},
		}),
	)
}
