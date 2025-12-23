// 03-params: 参数获取
//
// 📌 Gin 参数获取方式:
//   - c.Param("id")       - 路径参数 /users/:id
//   - c.Query("q")        - 查询参数 ?q=xxx
//   - c.DefaultQuery()    - 带默认值的查询参数
//   - c.PostForm()        - 表单参数
//   - c.ShouldBindJSON()  - JSON 绑定
//
// 📌 最佳实践:
//   - 路径参数用于资源标识
//   - 查询参数用于过滤/分页
//   - 请求体用于创建/更新数据
package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 路径参数
	r.GET("/users/:id", getUser)
	r.GET("/users/:userId/orders/:orderId", getOrder)

	// 查询参数
	r.GET("/search", search)

	// 表单参数
	r.POST("/login-form", loginForm)

	// 混合参数
	r.PUT("/users/:id", updateUser)

	r.Run(":8080")
}

// 路径参数示例
// GET /users/123
func getUser(c *gin.Context) {
	id := c.Param("id")

	// 转换为整数
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户 ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":     idInt,
		"user_id_str": id,
	})
}

// 多路径参数
// GET /users/1/orders/100
func getOrder(c *gin.Context) {
	userId := c.Param("userId")
	orderId := c.Param("orderId")

	c.JSON(http.StatusOK, gin.H{
		"user_id":  userId,
		"order_id": orderId,
	})
}

// 查询参数示例
// GET /search?q=golang&page=2&limit=10&tags=web,api
func search(c *gin.Context) {
	// 必填参数
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "q 参数必填"})
		return
	}

	// 带默认值的参数
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	// 获取数组参数方式1: ?tags=web,api
	tags := c.Query("tags")

	// 获取数组参数方式2: ?tag=web&tag=api
	tagArray := c.QueryArray("tag")

	c.JSON(http.StatusOK, gin.H{
		"query":     q,
		"page":      page,
		"limit":     limit,
		"tags":      tags,
		"tag_array": tagArray,
	})
}

// 表单参数示例
// POST /login-form
// Content-Type: application/x-www-form-urlencoded
func loginForm(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	remember := c.DefaultPostForm("remember", "false")

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名和密码必填"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"remember": remember,
		"message":  "登录成功",
	})
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// 混合参数：路径 + 请求体
// PUT /users/123
func updateUser(c *gin.Context) {
	// 路径参数
	id := c.Param("id")

	// JSON 请求体
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       id,
		"username": req.Username,
		"email":    req.Email,
		"message":  "更新成功",
	})
}
