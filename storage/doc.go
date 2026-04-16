// Package storage 提供对象存储抽象层。
//
// 通过统一的 [Store] 接口封装上传、下载、列举、删除等操作，
// 并以 [Manager] 作为线程安全的对外入口。
//
// 当前内置实现：
//   - provider/local —— 本地文件系统（包名 local）
//   - provider/aliyunoss —— 阿里云 OSS（包名 aliyunoss）
//
// 扩展其他后端只需实现 [Store] 接口。
package storage
