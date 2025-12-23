// Package utils 提供通用工具函数
//
// 📌 pkg 目录特性:
//   - 可以被外部模块导入使用
//   - 适合放通用、可复用的代码
//   - 不应包含业务逻辑
package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateID 生成随机 ID
func GenerateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GenerateToken 生成随机 Token
func GenerateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
