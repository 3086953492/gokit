# YaBase Logger 使用指南

YaBase Logger 是一个基于 Zap 的高性能日志库，提供了灵活的配置选项、日志轮转、自动清理等功能。

## 功能特性

- 🚀 基于 Uber Zap 的高性能日志记录
- 📁 支持文件和控制台双重输出
- 🔄 支持按大小和按日期两种日志轮转方式
- 🗂️ 自动日志文件清理和压缩
- ⚙️ 灵活的配置选项
- 🔧 链式配置 API
- 🛡️ 线程安全

## 快速开始

### 1. 基本使用

```go
package main

import (
    "github.com/3086953492/YaBase/logger"
    "go.uber.org/zap"
)

func main() {
    // 使用默认配置
    logger.Info("应用启动", zap.String("version", "1.0.0"))
    logger.Debug("调试信息", zap.Int("count", 42))
    logger.Warn("警告信息", zap.String("component", "database"))
    logger.Error("错误信息", zap.Error(err))
}
```

### 2. 使用构建器自定义配置

```go
package main

import (
    "github.com/3086953492/YaBase/logger"
    "go.uber.org/zap"
)

func main() {
    // 创建自定义logger
    customLogger, err := logger.NewBuilder().
        WithLevel("debug").
        WithFilename("logs/myapp.log").
        WithConsole(true).
        WithRotateDaily(true).
        WithRotationConfig(50, 5, 30, true).
        Build()
    
    if err != nil {
        panic(err)
    }
    
    // 设置为默认logger
    logger.SetDefault(customLogger)
    
    // 使用logger
    logger.Info("自定义配置的logger已启动")
}
```

### 3. 使用配置结构体初始化

```go
package main

import (
    "github.com/3086953492/YaBase/configs"
    "github.com/3086953492/YaBase/logger"
)

func main() {
    // 创建配置
    config := configs.LogConfig{
        Level:       "info",
        Filename:    "logs/app.log",
        MaxSize:     100,    // 100MB
        MaxBackups:  5,      // 保留5个备份
        MaxAge:      30,     // 保留30天
        Compress:    true,   // 压缩旧文件
        RotateDaily: true,   // 按日期轮转
        Console:     true,   // 输出到控制台
        LogsDir:     "logs", // 日志目录
    }
    
    // 使用配置初始化
    err := logger.InitWithConfig(config)
    if err != nil {
        panic(err)
    }
    
    logger.Info("使用配置结构体初始化的logger")
}
```

## 详细配置选项

### LogConfig 结构体

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| Level | string | "info" | 日志级别 (debug/info/warn/error) |
| Filename | string | "logs/app.log" | 日志文件路径 |
| MaxSize | int | 100 | 单个文件最大大小(MB) |
| MaxBackups | int | 3 | 最大备份文件数 |
| MaxAge | int | 7 | 文件最大保存天数 |
| Compress | bool | true | 是否压缩旧文件 |
| RotateDaily | bool | true | 是否按日期轮转 |
| Console | bool | true | 是否同时输出到控制台 |
| LogsDir | string | "logs" | 日志目录 |

### 构建器方法

```go
// 创建构建器
builder := logger.NewBuilder()

// 设置完整配置
builder.WithConfig(config)

// 设置日志级别
builder.WithLevel("debug")

// 设置日志文件名
builder.WithFilename("logs/app.log")

// 设置是否按日期轮转
builder.WithRotateDaily(true)

// 设置是否输出到控制台
builder.WithConsole(true)

// 设置轮转配置
builder.WithRotationConfig(maxSize, maxBackups, maxAge, compress)

// 构建logger
logger, err := builder.Build()
```

## 使用示例

### 1. 基础日志记录

```go
package main

import (
    "github.com/3086953492/YaBase/logger"
    "go.uber.org/zap"
)

func main() {
    // 不同级别的日志
    logger.Debug("调试信息", zap.String("module", "auth"))
    logger.Info("用户登录", zap.String("username", "john"), zap.String("ip", "192.168.1.1"))
    logger.Warn("数据库连接缓慢", zap.Duration("duration", time.Second*2))
    logger.Error("数据库连接失败", zap.Error(err))
    
    // 结构化日志
    logger.Info("处理请求",
        zap.String("method", "POST"),
        zap.String("path", "/api/users"),
        zap.Int("status", 200),
        zap.Duration("duration", time.Millisecond*150),
    )
}
```

### 2. 错误日志记录

**注意：** 推荐使用 `errors` 包的简化 API 进行错误日志记录，它会自动获取函数名并记录日志。

```go
package main

import (
    "github.com/3086953492/YaBase/errors"
)

func handleUser(userID string) (*User, error) {
    user, err := repository.GetUser(userID)
    if err != nil {
        // 使用 errors 包的 Log() 方法，自动记录日志
        return nil, errors.Database().
            Msg("用户验证失败").
            Err(err).
            Field("user_id", userID).
            Field("action", "validate").
            Log()
    }
    return user, nil
}
```

如果需要直接使用 logger 包记录非错误日志：

```go
package main

import (
    "github.com/3086953492/YaBase/logger"
    "go.uber.org/zap"
)

func handleUser(userID string) {
    logger.Info("处理用户请求",
        zap.String("user_id", userID),
        zap.String("action", "validate"),
    )
}
```

### 3. 获取原生 Zap Logger

```go
package main

import (
    "github.com/3086953492/YaBase/logger"
    "go.uber.org/zap"
)

func main() {
    // 获取原生zap logger进行高级操作
    zapLogger := logger.GetLogger()
    
    // 使用zap的高级功能
    zapLogger.With(
        zap.String("service", "user-service"),
        zap.String("version", "1.0.0"),
    ).Info("服务启动")
}
```

### 4. 自定义 Logger 实例

```go
package main

import (
    "github.com/3086953492/YaBase/logger"
    "go.uber.org/zap"
)

func main() {
    // 创建自定义logger实例
    customLogger := &logger.Logger{}
    
    // 设置zap logger
    zapLogger, _ := logger.NewBuilder().
        WithLevel("debug").
        WithFilename("logs/custom.log").
        Build()
    
    customLogger.SetLogger(zapLogger)
    
    // 使用自定义实例
    customLogger.Info("自定义logger实例")
}
```

## 日志轮转

### 按大小轮转

```go
// 当文件达到100MB时轮转，保留5个备份，保留30天，压缩旧文件
logger, err := logger.NewBuilder().
    WithFilename("logs/app.log").
    WithRotateDaily(false). // 关闭按日期轮转
    WithRotationConfig(100, 5, 30, true).
    Build()
```

### 按日期轮转

```go
// 每天轮转，文件名格式：app-2024-01-15.log
logger, err := logger.NewBuilder().
    WithFilename("logs/app.log").
    WithRotateDaily(true). // 启用按日期轮转
    WithRotationConfig(100, 5, 30, true).
    Build()
```

## 日志清理

系统会自动清理过期的日志文件：

- 每天凌晨1点执行清理任务
- 删除超过 `MaxAge` 天数的文件
- 保留不超过 `MaxBackups` 个备份文件
- 自动压缩旧文件（如果启用）

## 性能优化建议

1. **合理设置日志级别**：生产环境建议使用 `info` 或 `warn` 级别
2. **启用压缩**：减少磁盘空间占用
3. **合理设置轮转参数**：避免日志文件过大
4. **使用结构化日志**：便于后续分析和查询

## 注意事项

1. 确保日志目录有写入权限
2. 定期检查日志文件大小，避免磁盘空间不足
3. 在生产环境中建议关闭 `Console` 输出以提高性能
4. 使用 `zap.Field` 进行结构化日志记录，避免字符串拼接

## 依赖项

- `go.uber.org/zap` - 核心日志库
- `gopkg.in/natefinch/lumberjack.v2` - 日志轮转库

## 许可证

本项目采用 MIT 许可证。
