// Package model 定义数据模型（最底层，被所有层依赖）
package model

import "time"

// User 用户模型
type User struct {
	ID        string
	Username  string
	Email     string
	CreatedAt time.Time
}
