# Response 包使用指南

## 概述

`response` 包为 gokit 项目提供了统一的 HTTP 响应处理机制，支持成功响应、错误响应、分页响应和重定向。该包通过依赖注入的方式解耦配置，支持灵活的自定义策略。

## 核心特性

- **统一响应格式**：所有 API 返回一致的 JSON 结构
- **灵活配置**：通过 `Init` 函数注入运行环境和策略
- **解耦设计**：不直接依赖全局配置，便于测试和复用
- **错误映射**：自动将 `AppError` 映射为相应的 HTTP 状态码
- **OAuth 支持**：内置 OAuth 2.0 错误码映射
- **重定向支持**：支持临时和永久重定向，可携带错误信息

## 初始化配置

在应用启动时（通常在 `main.go` 中），需要先初始化 `response` 包：

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/3086953492/gokit/response"
    "github.com/3086953492/gokit/config"
)

func main() {
    // 加载配置
    cfg := config.GetGlobalConfig()
    
    // 初始化 response 包
    response.Init(
        // 根据运行环境决定是否显示错误详情
        response.WithShowErrorDetail(cfg.Server.Mode == gin.DebugMode),
        
        // 可选：自定义默认消息
        response.WithDefaultSuccessMessage("操作成功"),
        response.WithDefaultPaginatedMessage("获取成功"),
        
        // 可选：自定义错误映射（如果需要）
        // response.WithErrorStatusMapper(customMapper),
        // response.WithOAuthErrorCodeMapper(customOAuthMapper),
    )
    
    // 启动服务器
    // ...
}
```

### 配置选项

| 选项 | 说明 | 默认值 |
|-----|------|-------|
| `WithShowErrorDetail(bool)` | 是否显示错误详细信息（Cause、Fields） | `false` |
| `WithErrorStatusMapper(func)` | 自定义错误类型到 HTTP 状态码的映射 | 内置默认映射 |
| `WithOAuthErrorCodeMapper(func)` | 自定义错误类型到 OAuth 2.0 错误码的映射 | 内置默认映射 |
| `WithDefaultSuccessMessage(string)` | 默认成功消息 | `"操作成功"` |
| `WithDefaultPaginatedMessage(string)` | 默认分页成功消息 | `"获取成功"` |
| `WithFallbackErrorMessage(string)` | 兜底错误消息（非 AppError） | `"系统内部错误"` |
| `WithFallbackErrorCode(string)` | 兜底错误代码（非 AppError） | `"SYSTEM_ERROR"` |

## 核心结构

```go
type Response struct {
    Success bool   `json:"success"`           // 请求是否成功
    Error   string `json:"error,omitempty"`   // 错误类型（仅错误时）
    Message string `json:"message"`           // 响应消息
    Data    any    `json:"data,omitempty"`    // 响应数据（可选）
}

type PaginatedData struct {
    Items      any   `json:"items"`        // 数据项列表
    Total      int64 `json:"total"`        // 总记录数
    Page       int   `json:"page"`         // 当前页码
    PageSize   int   `json:"page_size"`    // 每页大小
    TotalPages int64 `json:"total_pages"`  // 总页数
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

// message 为空时使用配置的默认消息
response.Success(c, "", nil)

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
// 使用 AppError
if user == nil {
    response.Error(c, errors.NotFound().Msg("用户不存在").Build())
    return
}

// 使用预定义错误
response.Error(c, errors.ErrUserNotFound)

// 包装数据库错误
if err := db.Create(&user).Error; err != nil {
    response.Error(c, errors.Internal().Msg("创建用户失败").Err(err).Build())
    return
}

// 普通 error 会被自动包装为 500 错误
if err := someOperation(); err != nil {
    response.Error(c, err)
    return
}
```

**默认错误类型与 HTTP 状态码映射：**

| 错误类型 | HTTP 状态码 | 说明 |
|---------|-----------|------|
| `NOT_FOUND` | 404 | 资源不存在 |
| `INVALID_INPUT` | 400 | 输入参数错误 |
| `UNAUTHORIZED` | 401 | 未授权 |
| `FORBIDDEN` | 403 | 权限不足 |
| `DUPLICATE` | 409 | 数据重复 |
| `VALIDATION` | 422 | 数据验证失败 |
| 其他 | 500 | 内部错误 |

**响应格式（生产环境）：**
```json
{
    "success": false,
    "error": "NOT_FOUND",
    "message": "用户不存在"
}
```

**响应格式（开发环境，ShowErrorDetail=true）：**
```json
{
    "success": false,
    "error": "错误消息: 用户不存在 , 错误类型: NOT_FOUND , 错误详情: sql: no rows in result set , 错误字段: map[user_id:123]",
    "message": "用户不存在"
}
```

### 3. 分页响应 (Paginated)

用于返回分页数据，包含分页元信息。总页数会自动向上取整。

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

### 4. 重定向 (Redirect)

支持临时重定向（302）和永久重定向（301），可携带自定义参数和错误信息。

```go
func RedirectTemporary(c *gin.Context, targetURL string, err error, params map[string]string)
func RedirectPermanent(c *gin.Context, targetURL string, err error, params map[string]string)
```

**参数优先级（从低到高）：**
1. 原 URL 中的查询参数
2. `params` 中的自定义参数（覆盖原 URL 中的同名参数）
3. 错误相关参数 `error` 和 `error_description`（覆盖 `params` 中的同名参数）

**使用示例：**

```go
// 基础重定向
response.RedirectTemporary(c, "https://example.com/callback", nil, nil)

// 带自定义参数的重定向
response.RedirectTemporary(c, "https://example.com/callback", nil, map[string]string{
    "code": "auth_code_123",
    "state": "random_state",
})

// 带错误信息的重定向（OAuth 2.0 风格）
err := errors.Unauthorized().Msg("用户未授权").Build()
response.RedirectTemporary(c, "https://example.com/callback", err, map[string]string{
    "state": "random_state",
})
// 重定向到: https://example.com/callback?state=random_state&error=unauthorized_client&error_description=用户未授权

// 永久重定向
response.RedirectPermanent(c, "https://example.com/new-location", nil, nil)
```

**OAuth 2.0 错误码映射：**

| AppError 类型 | OAuth 错误码 |
|--------------|-------------|
| `INVALID_INPUT` | `invalid_request` |
| `UNAUTHORIZED` | `unauthorized_client` |
| `FORBIDDEN` | `access_denied` |
| 其他 | `server_error` |

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
        response.Error(c, errors.InvalidInput().Msg("无效的用户ID").Build())
        return
    }
    
    var user User
    if err := db.First(&user, userID).Error; err != nil {
        response.Error(c, errors.NotFound().Msg("用户不存在").Err(err).Build())
        return
    }
    
    response.Success(c, "用户信息获取成功", user)
}

// 创建用户
func CreateUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        response.Error(c, errors.Validation().Msg("请求数据格式错误").Err(err).Build())
        return
    }
    
    // 验证用户名是否已存在
    var existingUser User
    if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
        response.Error(c, errors.Duplicate().Msg("用户名已存在").Build())
        return
    }
    
    if err := db.Create(&user).Error; err != nil {
        response.Error(c, errors.Internal().Msg("创建用户失败").Err(err).Build())
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
        response.Error(c, errors.Internal().Msg("查询用户列表失败").Err(err).Build())
        return
    }
    
    response.Paginated(c, users, total, page, pageSize)
}

// OAuth 授权回调示例
func OAuthCallback(c *gin.Context) {
    code := c.Query("code")
    if code == "" {
        err := errors.InvalidInput().Msg("缺少授权码").Build()
        response.RedirectTemporary(c, "https://app.example.com/login", err, nil)
        return
    }
    
    // 验证授权码
    token, err := validateAuthCode(code)
    if err != nil {
        appErr := errors.Unauthorized().Msg("授权码无效").Err(err).Build()
        response.RedirectTemporary(c, "https://app.example.com/login", appErr, nil)
        return
    }
    
    // 成功，重定向到应用
    response.RedirectTemporary(c, "https://app.example.com/dashboard", nil, map[string]string{
        "token": token,
    })
}
```

## 最佳实践

### 1. 初始化配置
- 在应用启动时调用 `response.Init()` 一次
- 根据运行环境（开发/生产）决定是否显示错误详情
- 使用有意义的默认消息

### 2. 错误处理
- 优先使用 `errors` 包的 Builder 模式创建错误
- 使用 `.Err()` 包装底层错误，便于调试
- 使用 `.Field()` 添加上下文信息
- 在开发环境启用详细错误信息，生产环境关闭

### 3. 成功响应
- 提供有意义的成功消息
- 根据业务需求决定是否返回数据
- 对于列表数据，考虑使用分页响应

### 4. 分页响应
- 始终验证分页参数的有效性
- 设置合理的默认值和最大值限制
- 总页数会自动向上取整

### 5. 重定向
- 明确使用临时（302）还是永久（301）重定向
- 利用参数优先级机制，确保错误信息不被覆盖
- 对于 OAuth 场景，使用标准的错误码映射

### 6. 自定义映射
- 如果需要自定义 HTTP 状态码映射，使用 `WithErrorStatusMapper`
- 如果需要自定义 OAuth 错误码映射，使用 `WithOAuthErrorCodeMapper`
- 保持映射逻辑的一致性和可预测性

## 注意事项

1. **必须初始化**：虽然未初始化时会使用默认配置，但建议显式调用 `Init()`
2. **只初始化一次**：`Init()` 使用 `sync.Once` 保证只执行一次，多次调用无效
3. **线程安全**：配置读取是线程安全的，可以在并发环境中使用
4. **数据序列化**：确保返回的数据结构能够正确序列化为 JSON
5. **性能考虑**：对于大量数据的分页查询，注意数据库查询性能
6. **安全性**：生产环境不要显示详细错误信息，避免暴露敏感系统信息

## 架构设计

### 文件结构

```
response/
├── config.go      # 配置管理、初始化、Option 模式
├── success.go     # 成功响应、分页响应
├── error.go       # 错误响应、错误映射
├── redirect.go    # 重定向逻辑
└── README.md      # 使用文档
```

### 依赖关系

- `response` 包只依赖 `errors` 包（用于识别 `AppError`）和 `gin`
- 不直接依赖全局 `config` 包，通过 `Init()` 注入配置
- 各模块职责单一，易于维护和测试

通过使用 `response` 包，您可以确保 API 响应的一致性和规范性，提高代码的可维护性和用户体验。
