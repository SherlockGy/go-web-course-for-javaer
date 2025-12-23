// 05-errors: 错误处理模式
//
// 📌 最佳实践:
//   - error 作为最后一个返回值
//   - 调用后立即检查 if err != nil
//   - 使用 fmt.Errorf 包装错误，添加上下文
//   - 使用 errors.Is/As 判断错误类型
//   - panic 仅用于不可恢复的错误
//
// 🆚 Java 对比:
//   Java: try { ... } catch (Exception e) { ... }
//   Go:   result, err := doSomething(); if err != nil { ... }
//
//   Go 的方式更"显式"，强制你思考"如果失败了怎么办"
package main

import (
	"errors"
	"fmt"
	"os"
)

// 预定义错误 - 可用于 errors.Is 比较
var (
	ErrNotFound     = errors.New("资源不存在")
	ErrUnauthorized = errors.New("未授权")
	ErrInvalidInput = errors.New("输入无效")
)

// 自定义错误类型 - 携带更多信息
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("验证失败 [%s]: %s", e.Field, e.Message)
}

func main() {
	// === 基本错误处理 ===
	result, err := divide(10, 0)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("结果: %.2f\n", result)
	}

	// === 错误包装（Go 1.13+）===
	err = readConfig("config.yaml")
	if err != nil {
		fmt.Printf("配置错误: %v\n", err)

		// 检查是否包含特定错误
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("-> 文件不存在，使用默认配置")
		}
	}

	// === 预定义错误 ===
	user, err := findUser(999)
	if err != nil {
		// 使用 errors.Is 比较
		if errors.Is(err, ErrNotFound) {
			fmt.Println("用户不存在")
		} else {
			fmt.Printf("查找错误: %v\n", err)
		}
	} else {
		fmt.Printf("用户: %s\n", user)
	}

	// === 自定义错误类型 ===
	err = validateUser("", "invalid-email")
	if err != nil {
		// 使用 errors.As 提取错误详情
		var validationErr *ValidationError
		if errors.As(err, &validationErr) {
			fmt.Printf("字段: %s, 原因: %s\n", validationErr.Field, validationErr.Message)
		}
	}

	// === 多错误处理（Go 1.20+）===
	errs := validateAll("", "bad", "")
	if errs != nil {
		fmt.Printf("多个错误: %v\n", errs)
	}

	// === panic 和 recover ===
	safeDivide(10, 0)
	fmt.Println("程序继续运行...")
}

// 基本错误返回
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("除数不能为零")
	}
	return a / b, nil
}

// 错误包装 - 添加上下文
func readConfig(filename string) error {
	_, err := os.ReadFile(filename)
	if err != nil {
		// 使用 %w 包装原始错误（保留错误链）
		return fmt.Errorf("读取配置文件 %s 失败: %w", filename, err)
	}
	return nil
}

// 使用预定义错误
func findUser(id int) (string, error) {
	users := map[int]string{1: "Tom", 2: "Jerry"}

	if user, ok := users[id]; ok {
		return user, nil
	}
	// 包装预定义错误并添加上下文
	return "", fmt.Errorf("查找用户 id=%d: %w", id, ErrNotFound)
}

// 返回自定义错误类型
func validateUser(username, email string) error {
	if username == "" {
		return &ValidationError{Field: "username", Message: "不能为空"}
	}
	if email == "" || len(email) < 5 {
		return &ValidationError{Field: "email", Message: "格式无效"}
	}
	return nil
}

// 多错误合并（Go 1.20+）
func validateAll(name, email, phone string) error {
	var errs []error

	if name == "" {
		errs = append(errs, fmt.Errorf("name: %w", ErrInvalidInput))
	}
	if email == "" {
		errs = append(errs, fmt.Errorf("email: %w", ErrInvalidInput))
	}
	if phone == "" {
		errs = append(errs, fmt.Errorf("phone: %w", ErrInvalidInput))
	}

	if len(errs) > 0 {
		return errors.Join(errs...) // Go 1.20+
	}
	return nil
}

// panic 和 recover - 仅用于不可恢复的错误
func safeDivide(a, b int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获 panic: %v\n", r)
		}
	}()

	if b == 0 {
		panic("除数为零！") // 不推荐用于普通错误
	}
	fmt.Printf("%d / %d = %d\n", a, b, a/b)
}
