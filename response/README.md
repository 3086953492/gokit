# Response 包使用指南

## 概述

`response` 包为 gokit 项目提供了统一的 HTTP 响应处理机制，支持成功响应、错误响应和分页响应。该包与 `errors` 包紧密集成，能够自动处理不同类型的应用错误并返回相应的 HTTP 状态码。

## 核心结构

```go
type Response struct {
    Success bool   `json:"success"`           // 请求是否成功
    Error   string `json:"error,omitempty"`   // 错误类型（仅错误时）
    Message string `json:"message"`           // 响应消息
    Data    any    `json:"data,omitempty"`    // 响应数据（可选）
}
```

## 主要功能

### 1. 成功响应 (Success)

用于返回成功的业务操作结果。

```go
func Success(c *gin.Context, message string, data any)
```

**使用示例：**

```go
// 基础用法 - 只返回消息
response.Success(c, "操作成功", nil)

// 返回数据
user := User{
    ID:   1,
    Name: "张三",
    Email: "zhangsan@example.com",
}
response.Success(c, "用户信息获取成功", user)

// 返回列表数据
users := []User{...}
response.Success(c, "用户列表获取成功", users)
```

**响应格式：**
```json
{
    "success": true,
    "message": "操作成功",
    "data": {
        "id": 1,
        "name": "张三",
        "email": "zhangsan@example.com"
    }
}
```

### 2. 错误响应 (Error)

用于返回各种类型的错误信息，自动根据错误类型设置相应的 HTTP 状态码。

```go
func Error(c *gin.Context, err error)
```

**使用示例：**

```go
// 使用预定义错误
if user == nil {
    response.Error(c, errors.ErrUserNotFound)
    return
}

// 使用自定义错误
if err := validateUser(user); err != nil {
    response.Error(c, errors.New(errors.TypeValidation, "用户数据验证失败"))
    return
}

// 包装数据库错误
if err := db.Create(&user).Error; err != nil {
    response.Error(c, errors.FromDatabaseError(err))
    return
}
```

**错误类型与HTTP状态码映射：**

| 错误类型 | HTTP状态码 | 说明 |
|---------|-----------|------|
| `UNAUTHORIZED` | 401 | 未授权 |
| `TOKEN_EXPIRED` | 401 | 令牌过期 |
| `TOKEN_INVALID` | 401 | 无效令牌 |
| `PERMISSION_DENIED` | 403 | 权限不足 |
| `NOT_FOUND` | 404 | 资源不存在 |
| `USER_NOT_FOUND` | 404 | 用户不存在 |
| `INVALID_INPUT` | 400 | 输入参数错误 |
| `VALIDATION_ERROR` | 400 | 数据验证失败 |
| `DUPLICATE_KEY` | 400 | 数据重复 |
| `INTERNAL_ERROR` | 500 | 系统内部错误 |
| `DATABASE_ERROR` | 500 | 数据库错误 |

**响应格式：**
```json
{
    "success": false,
    "error": "USER_NOT_FOUND",
    "message": "用户不存在"
}
```

### 3. 分页响应 (Paginated)

用于返回分页数据，包含分页元信息。

```go
func Paginated(c *gin.Context, data any, total int64, page, pageSize int)
```

**使用示例：**

```go
// 获取用户列表（分页）
func GetUsers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    
    var users []User
    var total int64
    
    // 查询数据
    db.Model(&User{}).Count(&total)
    db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)
    
    // 返回分页响应
    response.Paginated(c, users, total, page, pageSize)
}
```

**响应格式：**
```json
{
    "success": true,
    "message": "获取成功",
    "data": {
        "items": [
            {"id": 1, "name": "张三"},
            {"id": 2, "name": "李四"}
        ],
        "total": 100,
        "page": 1,
        "page_size": 10,
        "total_pages": 10
    }
}
```

## 完整使用示例

### 用户管理接口示例

```go
package handlers

import (
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/3086953492/gokit/response"
    "github.com/3086953492/gokit/errors"
)

// 获取用户信息
func GetUser(c *gin.Context) {
    userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        response.Error(c, errors.New(errors.TypeInvalidInput, "无效的用户ID"))
        return
    }
    
    var user User
    if err := db.First(&user, userID).Error; err != nil {
        response.Error(c, errors.FromDatabaseError(err))
        return
    }
    
    response.Success(c, "用户信息获取成功", user)
}

// 创建用户
func CreateUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        response.Error(c, errors.New(errors.TypeValidation, "请求数据格式错误"))
        return
    }
    
    // 验证用户名是否已存在
    var existingUser User
    if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
        response.Error(c, errors.ErrUsernameExists)
        return
    }
    
    if err := db.Create(&user).Error; err != nil {
        response.Error(c, errors.FromDatabaseError(err))
        return
    }
    
    response.Success(c, "用户创建成功", user)
}

// 获取用户列表（分页）
func GetUsers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 10
    }
    
    var users []User
    var total int64
    
    // 构建查询
    query := db.Model(&User{})
    
    // 添加搜索条件
    if name := c.Query("name"); name != "" {
        query = query.Where("name LIKE ?", "%"+name+"%")
    }
    
    // 获取总数
    query.Count(&total)
    
    // 分页查询
    if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
        response.Error(c, errors.FromDatabaseError(err))
        return
    }
    
    response.Paginated(c, users, total, page, pageSize)
}

// 更新用户
func UpdateUser(c *gin.Context) {
    userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        response.Error(c, errors.New(errors.TypeInvalidInput, "无效的用户ID"))
        return
    }
    
    var user User
    if err := db.First(&user, userID).Error; err != nil {
        response.Error(c, errors.FromDatabaseError(err))
        return
    }
    
    var updateData User
    if err := c.ShouldBindJSON(&updateData); err != nil {
        response.Error(c, errors.New(errors.TypeValidation, "请求数据格式错误"))
        return
    }
    
    if err := db.Model(&user).Updates(updateData).Error; err != nil {
        response.Error(c, errors.FromDatabaseError(err))
        return
    }
    
    response.Success(c, "用户更新成功", user)
}

// 删除用户
func DeleteUser(c *gin.Context) {
    userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        response.Error(c, errors.New(errors.TypeInvalidInput, "无效的用户ID"))
        return
    }
    
    if err := db.Delete(&User{}, userID).Error; err != nil {
        response.Error(c, errors.FromDatabaseError(err))
        return
    }
    
    response.Success(c, "用户删除成功", nil)
}
```

## 最佳实践

### 1. 错误处理
- 优先使用预定义的错误类型和实例
- 对于数据库错误，使用 `errors.FromDatabaseError()` 进行转换
- 对于验证错误，提供具体的错误信息

### 2. 成功响应
- 提供有意义的成功消息
- 根据业务需求决定是否返回数据
- 对于列表数据，考虑使用分页响应

### 3. 分页响应
- 始终验证分页参数的有效性
- 设置合理的默认值和最大值限制
- 提供完整的分页元信息

### 4. 响应一致性
- 所有接口都使用统一的响应格式
- 错误消息要清晰明确，便于前端处理
- 保持响应结构的一致性

## 注意事项

1. **错误类型映射**：确保 `errors` 包中定义的错误类型与 `getHTTPStatus` 函数中的映射保持一致
2. **数据序列化**：确保返回的数据结构能够正确序列化为 JSON
3. **性能考虑**：对于大量数据的分页查询，注意数据库查询性能
4. **安全性**：避免在错误响应中暴露敏感的系统信息

通过使用 `response` 包，您可以确保 API 响应的一致性和规范性，提高代码的可维护性和用户体验。
