// Package greeting 提供问候相关功能
//
// 📌 Go 文档注释规范:
//   - 包注释以 "Package xxx" 开头
//   - 函数注释以函数名开头
//   - 使用 go doc 命令查看文档
package greeting

import "fmt"

// MaxNameLength 是名字的最大长度（导出常量）
const MaxNameLength = 50

// DefaultLanguage 是默认语言（导出变量）
var DefaultLanguage = "zh-CN"

// Hello 返回问候语（导出函数 - 大写开头）
func Hello(name string) string {
	formatted := formatName(name) // 调用私有函数
	return fmt.Sprintf("你好, %s!", formatted)
}

// FormalGreeting 返回正式问候语
func FormalGreeting(title, name string) string {
	return fmt.Sprintf("尊敬的%s%s，您好！", title, name)
}

// formatName 格式化名字（私有函数 - 小写开头）
// 只能在 greeting 包内部调用
func formatName(name string) string {
	if len(name) > MaxNameLength {
		return name[:MaxNameLength]
	}
	return name
}

// internalHelper 是另一个私有函数
func internalHelper() {
	// 包内部使用，外部不可见
}
