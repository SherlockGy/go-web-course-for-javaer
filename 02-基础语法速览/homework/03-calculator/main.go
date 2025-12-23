// 作业3：安全计算器
//
// 📌 学习目标：
//   - 理解 Go 的多返回值设计
//   - 使用 error 接口处理错误
//   - 定义哨兵错误（sentinel error）
//   - 使用 errors.Is 检查错误
//   - 使用 fmt.Errorf 包装错误
//
// 📌 要求：
//   1. 定义三个哨兵错误：
//      - ErrDivideByZero：除数为零
//      - ErrNegativeSqrt：负数开方
//      - ErrOverflow：整数溢出（可选，加分项）
//   2. 实现 SafeDivide(a, b int) (int, error)
//   3. 实现 SafeSqrt(n float64) (float64, error)
//   4. 实现 Calculate(op string, a, b float64) (float64, error)
//      - 支持 "+", "-", "*", "/"
//      - 未知操作符返回错误
//   5. 使用 errors.Is 检查具体错误类型
//
// 📌 提示：
//   - var ErrXxx = errors.New("错误描述")
//   - return 0, ErrDivideByZero
//   - return result, nil  // 成功时 error 为 nil
//   - math.Sqrt(n) 计算平方根
//   - errors.Is(err, ErrDivideByZero) 检查错误
//
// 🆚 与 Java 对比：
//   - Java: throw new ArithmeticException("除数不能为零")
//   - Go: return 0, ErrDivideByZero
//
// 📌 运行：go run main.go
package main

import (
	"errors"
	"fmt"
	"math"
)

// TODO: 1. 定义哨兵错误
// var ErrDivideByZero = errors.New("...")
// var ErrNegativeSqrt = errors.New("...")

// TODO: 2. 实现 SafeDivide
// func SafeDivide(a, b int) (int, error) {
//     ...
// }

// TODO: 3. 实现 SafeSqrt
// func SafeSqrt(n float64) (float64, error) {
//     ...
// }

// TODO: 4. 实现 Calculate（支持 +, -, *, /）
// 提示：对于除法，需要检查除数是否为零
// func Calculate(op string, a, b float64) (float64, error) {
//     switch op {
//     case "+":
//         ...
//     case "/":
//         // 注意检查 b == 0
//         ...
//     default:
//         return 0, fmt.Errorf("未知操作符: %s", op)
//     }
// }

func main() {
	fmt.Println("=== 安全除法 ===")
	// TODO: 5. 测试 SafeDivide
	// result, err := SafeDivide(10, 2)
	// if err != nil {
	//     fmt.Printf("错误: %v\n", err)
	// } else {
	//     fmt.Printf("10 / 2 = %d\n", result)
	// }

	// TODO: 6. 测试除零错误
	// result, err = SafeDivide(10, 0)
	// if errors.Is(err, ErrDivideByZero) {
	//     fmt.Println("捕获到除零错误")
	// }

	fmt.Println("\n=== 安全开方 ===")
	// TODO: 7. 测试 SafeSqrt
	// sqrt, err := SafeSqrt(16)
	// sqrt, err = SafeSqrt(-1)  // 测试负数

	fmt.Println("\n=== 计算器 ===")
	// TODO: 8. 测试 Calculate
	// res, err := Calculate("+", 10, 5)
	// res, err = Calculate("/", 10, 0)
	// res, err = Calculate("%", 10, 5)  // 未知操作符

	// 以下是占位代码，完成后删除
	fmt.Println("请完成作业")
	_ = errors.New
	_ = math.Sqrt
}
