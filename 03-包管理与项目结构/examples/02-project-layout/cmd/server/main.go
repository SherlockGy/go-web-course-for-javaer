// 02-project-layout: 标准项目布局示例
//
// 📌 标准目录结构:
//   cmd/        - 可执行文件入口，每个子目录一个程序
//   internal/   - 私有代码，不能被外部导入
//   pkg/        - 可复用的公共库
//
// 运行: go run ./cmd/server
package main

import (
	"fmt"
	"log"

	"project-layout/internal/service"
	"project-layout/pkg/utils"
)

func main() {
	// 使用 pkg 中的公共工具
	id := utils.GenerateID()
	fmt.Printf("生成的 ID: %s\n", id)

	// 使用 internal 中的业务逻辑
	userSvc := service.NewUserService()

	user, err := userSvc.CreateUser("tom", "tom@example.com")
	if err != nil {
		log.Fatalf("创建用户失败: %v", err)
	}

	fmt.Printf("创建用户成功: %+v\n", user)

	// 获取用户
	found, err := userSvc.GetUser(user.ID)
	if err != nil {
		log.Fatalf("获取用户失败: %v", err)
	}

	fmt.Printf("找到用户: %s <%s>\n", found.Username, found.Email)
}
