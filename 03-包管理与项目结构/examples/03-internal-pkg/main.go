// 03-internal-pkg: internal 与 pkg 的区别演示
package main

import (
	"fmt"

	"internal-pkg-demo/internal/secret"
	"internal-pkg-demo/pkg/public"
)

func main() {
	// 可以导入 internal 包（同一模块内）
	secretValue := secret.GetSecret()
	fmt.Printf("内部秘密: %s\n", secretValue)

	// 可以导入 pkg 包
	publicValue := public.GetPublicInfo()
	fmt.Printf("公开信息: %s\n", publicValue)

	// 注意：如果其他模块尝试导入 internal/secret
	// 会得到编译错误: use of internal package not allowed
}
