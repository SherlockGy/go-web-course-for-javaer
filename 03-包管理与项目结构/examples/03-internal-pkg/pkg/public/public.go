// Package public 是可公开使用的工具包
//
// 📌 pkg 目录约定:
//   - 放置可被外部导入的代码
//   - 通常是通用工具、不含业务逻辑
//   - 注意：这只是约定，Go 编译器不会阻止导入
package public

// GetPublicInfo 返回公开信息
func GetPublicInfo() string {
	return "这是公开信息，任何人都可以访问"
}

// FormatMessage 格式化消息（通用工具）
func FormatMessage(msg string) string {
	return "[INFO] " + msg
}
