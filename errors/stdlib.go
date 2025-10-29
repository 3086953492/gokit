package errors

import stderrors "errors"

// 以下函数重新导出自 Go 标准库 errors 包
// 这样项目中只需要 import "github.com/3086953492/gokit/errors" 一个包即可

// NewSimple 创建一个简单的错误（标准库函数）
// 用于创建不需要结构化信息的简单错误
var NewSimple = stderrors.New

// As 从错误链中提取特定类型的错误（标准库函数）
// 用于类型断言，检查错误链中是否包含指定类型的错误
var As = stderrors.As

// Is 检查错误链中是否包含目标错误（标准库函数）
// 用于错误比较，支持错误链的递归检查
var Is = stderrors.Is

// Unwrap 解包错误，返回被包装的错误（标准库函数）
// 用于获取错误链中的下一个错误
var Unwrap = stderrors.Unwrap

// Join 合并多个错误为一个错误（标准库函数，Go 1.20+）
// 用于将多个错误组合成一个错误
var Join = stderrors.Join
