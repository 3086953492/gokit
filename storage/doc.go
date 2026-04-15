// Package storage 提供对象存储抽象层。
//
// 通过统一的 [Store] 接口封装上传、下载、列举、删除等操作，
// 并以 [Manager] 作为线程安全的对外入口。
//
// 当前内置实现：
//   - providerlocal —— 本地文件系统
//   - provideraliyunoss —— 阿里云 OSS
//
// 扩展其他后端只需实现 [Store] 接口，
// 若需支持 [Manager.DeleteByURL]，还需实现 [URLKeyResolver]。
package storage
