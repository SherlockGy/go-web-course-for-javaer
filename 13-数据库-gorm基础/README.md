# 13 - 数据库 GORM 基础

## 学习目标

掌握 GORM ORM 框架的基本使用。

---

## 🆚 Java 对比：ORM 设计哲学

| 特性 | JPA/Hibernate | GORM |
|------|---------------|------|
| 实体注解 | `@Entity`, `@Table` | struct tag |
| 主键 | `@Id @GeneratedValue` | `gorm.Model` 或 tag |
| 字段映射 | `@Column(name="...")` | `gorm:"column:..."` |
| 关联 | `@OneToMany`, `@ManyToOne` | `gorm:"foreignKey:..."` |
| 查询 | JPQL / Criteria API | 链式调用 |

> **洞察**：JPA 是"规范"，Hibernate 是实现，配置复杂但功能全面。GORM 是"约定优于配置"，简单直接，但高级功能需要原生 SQL。

---

## 学习要点

### 1. 安装与连接

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
```

```go
import (
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
)

func main() {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("连接数据库失败")
    }
}
```

### 2. 模型定义

```go
// gorm.Model 包含 ID, CreatedAt, UpdatedAt, DeletedAt
type User struct {
    gorm.Model
    Username string `gorm:"uniqueIndex;size:50"`
    Email    string `gorm:"uniqueIndex"`
    Password string `gorm:"size:100"`
    Age      int    `gorm:"default:0"`
}
```

> **🆚 JPA 对比**
> ```java
> @Entity
> @Table(name = "users")
> public class User {
>     @Id
>     @GeneratedValue(strategy = GenerationType.IDENTITY)
>     private Long id;
>
>     @Column(unique = true, length = 50)
>     private String username;
> }
> ```
> Go 用 struct tag 替代 Java 注解，更紧凑，但 IDE 支持不如注解好。

### 3. 自动迁移

```go
// 自动创建/更新表结构
db.AutoMigrate(&User{})
```

> **洞察**：GORM 的 AutoMigrate 只会"添加"字段，不会"删除"或"修改"字段类型。生产环境建议用专门的迁移工具（如 golang-migrate）。

### 4. CRUD 操作

```go
// Create
user := User{Username: "tom", Email: "tom@example.com"}
db.Create(&user)  // user.ID 会被自动赋值

// Read
var user User
db.First(&user, 1)                     // 按主键
db.First(&user, "username = ?", "tom") // 按条件
db.Find(&users)                        // 查询多条

// Update
db.Model(&user).Update("Age", 20)
db.Model(&user).Updates(User{Age: 20, Email: "new@example.com"})

// Delete
db.Delete(&user, 1)
```

> **🆚 JPA 对比**
> ```java
> // JPA: 先查后改
> User user = repo.findById(1L).get();
> user.setAge(20);
> repo.save(user);
>
> // GORM: 可以直接改
> db.Model(&User{}).Where("id = ?", 1).Update("Age", 20)
> ```
> GORM 更灵活，JPA 更"对象化"。

### 5. 常用 Tag

| Tag | 说明 | 示例 |
|-----|------|------|
| `primaryKey` | 主键 | `gorm:"primaryKey"` |
| `column` | 字段名 | `gorm:"column:user_name"` |
| `size` | 长度 | `gorm:"size:100"` |
| `unique` | 唯一 | `gorm:"unique"` |
| `uniqueIndex` | 唯一索引 | `gorm:"uniqueIndex"` |
| `index` | 索引 | `gorm:"index"` |
| `default` | 默认值 | `gorm:"default:0"` |
| `not null` | 非空 | `gorm:"not null"` |
| `-` | 忽略 | `gorm:"-"` |

---

## 示例代码

### examples/01-connect-db/
数据库连接

### examples/02-model-define/
模型定义与 Tag

### examples/03-basic-crud/
增删改查操作

---

## 作业任务

### 任务描述
定义 User 模型，实现完整的 CRUD 操作。

### 要求
1. User 模型字段：
   - `ID`（主键）
   - `Username`（唯一，最长 50）
   - `Email`（唯一）
   - `Password`
   - `CreatedAt`
   - `UpdatedAt`

2. 实现 UserRepository：
   - `Create(user *User) error`
   - `FindByID(id uint) (*User, error)`
   - `FindByUsername(username string) (*User, error)`
   - `Update(user *User) error`
   - `Delete(id uint) error`

### 验收标准
```go
repo := NewUserRepository(db)

// 创建
user := &User{Username: "tom", Email: "tom@example.com"}
repo.Create(user)

// 查询
found, _ := repo.FindByID(user.ID)
fmt.Println(found.Username)  // tom

// 更新
user.Email = "new@example.com"
repo.Update(user)

// 删除
repo.Delete(user.ID)
```

---

## 参考资料
- [GORM 官方文档](https://gorm.io/docs/)
