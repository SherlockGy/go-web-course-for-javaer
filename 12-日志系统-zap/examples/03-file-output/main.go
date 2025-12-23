// 03-file-output: 日志输出到文件
//
// 📌 生产环境日志策略:
//   - 同时输出到控制台和文件
//   - 文件按大小/时间轮转
//   - 保留最近 N 天的日志
//
// 📌 使用 lumberjack 进行日志轮转
package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	logger := createLoggerWithFile()
	defer logger.Sync()

	for i := 0; i < 10; i++ {
		logger.Info("测试日志",
			zap.Int("index", i),
			zap.String("message", "这是一条测试日志"),
		)
	}

	logger.Warn("警告日志")
	logger.Error("错误日志")
}

func createLoggerWithFile() *zap.Logger {
	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 文件输出（带轮转）
	fileWriter := &lumberjack.Logger{
		Filename:   "./logs/app.log", // 日志文件路径
		MaxSize:    10,               // 单个文件最大 10MB
		MaxBackups: 5,                // 保留最近 5 个备份
		MaxAge:     30,               // 保留最近 30 天
		Compress:   true,             // 压缩旧日志
	}

	// 控制台输出
	consoleWriter := zapcore.AddSync(os.Stdout)

	// 创建多输出核心
	core := zapcore.NewTee(
		// 文件输出：JSON 格式
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(fileWriter),
			zap.InfoLevel, // 文件只记录 Info 及以上
		),
		// 控制台输出：彩色格式
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			consoleWriter,
			zap.DebugLevel, // 控制台记录所有级别
		),
	)

	// 创建日志器
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))

	return logger
}
