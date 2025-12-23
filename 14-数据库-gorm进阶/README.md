# 14 - 数据库 GORM 进阶

## 学习目标

掌握 GORM 高级查询、关联关系和事务处理。

---

## 🆚 Java 对比：查询 API

```java
// JPA Criteria API（类型安全但繁琐）
CriteriaBuilder cb = em.getCriteriaBuilder();
CriteriaQuery<User> cq = cb.createQuery(User.class);
Root<User> root = cq.from(User.class);
cq.where(cb.equal(root.get("status"), 1));
```

```go
// GORM 链式调用（简洁直观）
db.Where("status = ?", 1).Find(&users)
```

> **洞察**：Go 选择了"简单"而非"类型安全"。写错字段名只有运行时才知道，但代码量少一半。

---

## 学习要点

### 1. 条件查询

```go
// Where
db.Where("age > ?", 18).Find(&users)
db.Where("name LIKE ?", "%tom%").Find(&users)
db.Where("name IN ?", []string{"tom", "jerry"}).Find(&users)

// Or
db.Where("age > ?", 18).Or("role = ?", "admin").Find(&users)

// Not
db.Not("status = ?", 0).Find(&users)

// 结构体查询（零值会被忽略）
db.Where(&User{Name: "tom", Age: 0}).Find(&users) // Age=0 不会作为条件
```

### 2. 分页查询

```go
// LIMIT + OFFSET
var users []User
db.Offset(10).Limit(10).Find(&users)  // 第 2 页，每页 10 条

// 分页封装
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        offset := (page - 1) * pageSize
        return db.Offset(offset).Limit(pageSize)
    }
}

db.Scopes(Paginate(2, 10)).Find(&users)
```

> **🆚 JPA 对比**
> ```java
> Pageable pageable = PageRequest.of(1, 10);
> Page<User> page = userRepo.findAll(pageable);
> ```
> Spring Data JPA 的分页更"自动化"，GORM 需要手动计算 offset。

### 3. 排序

```go
db.Order("created_at DESC").Find(&users)
db.Order("age DESC, name ASC").Find(&users)
```

### 4. 关联关系

```go
// 一对多：User 有多个 Order
type User struct {
    gorm.Model
    Username string
    Orders   []Order  // has many
}

type Order struct {
    gorm.Model
    UserID uint    // 外键
    Amount float64
    User   User    // belongs to
}

// 创建关联
db.AutoMigrate(&User{}, &Order{})

// 预加载查询（N+1 问题解决方案）
var user User
db.Preload("Orders").First(&user, 1)
```

> **🆚 JPA 对比**
> ```java
> @OneToMany(mappedBy = "user", fetch = FetchType.LAZY)
> private List<Order> orders;
> ```
> JPA 用 `FetchType.LAZY/EAGER`，GORM 用 `Preload()` 显式加载。Go 的方式更明确，不会有"懒加载在事务外失效"的坑。

### 5. 事务处理

```go
// 方式 1：自动事务
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err  // 返回错误会自动回滚
    }
    if err := tx.Create(&order).Error; err != nil {
        return err
    }
    return nil  // 返回 nil 自动提交
})

// 方式 2：手动事务
tx := db.Begin()
if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}
tx.Commit()
```

> **🆚 Spring 对比**
> ```java
> @Transactional
> public void createOrder(User user, Order order) {
>     userRepo.save(user);
>     orderRepo.save(order);  // 异常自动回滚
> }
> ```
> Spring 的 `@Transactional` 是声明式的（AOP），Go 是命令式的。Go 没有魔法，代码即逻辑。

### 6. 软删除

```go
// gorm.Model 自带 DeletedAt 字段
type User struct {
    gorm.Model  // 包含 DeletedAt
    Name string
}

// 软删除（设置 DeletedAt）
db.Delete(&user)

// 查询（默认排除已删除）
db.Find(&users)

// 查询包含已删除
db.Unscoped().Find(&users)

// 永久删除
db.Unscoped().Delete(&user)
```

---

## 示例代码

### examples/01-advanced-query/
高级查询（Where、Or、分页）

### examples/02-associations/
关联关系（一对多、预加载）

### examples/03-transactions/
事务操作

### examples/04-soft-delete/
软删除

---

## 作业任务

### 任务描述
实现用户和订单的一对多关联查询。

### 模型
```go
type User struct {
    gorm.Model
    Username string
    Orders   []Order
}

type Order struct {
    gorm.Model
    UserID  uint
    Product string
    Amount  float64
}
```

### 要求
1. 实现分页查询用户 `ListUsers(page, pageSize int)`
2. 查询用户时预加载订单列表
3. 实现创建订单的事务（检查用户存在 → 创建订单）

### 验收标准
```go
// 分页查询，带订单
users := ListUsers(1, 10)
fmt.Println(users[0].Orders)  // 用户的订单列表

// 创建订单（事务）
CreateOrder(userID, "iPhone", 9999.0)
```

---

## 参考资料
- [GORM 高级查询](https://gorm.io/docs/advanced_query.html)
- [GORM 关联](https://gorm.io/docs/associations.html)
