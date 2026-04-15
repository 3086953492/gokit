// Package provideraliyunoss 实现基于阿里云 OSS 的 [storage.Store] 后端。
//
// 通过 [New] 创建实例，配置项见 [Config]。
// 同时实现了 [storage.URLKeyResolver]，
// 支持 [storage.Manager.DeleteByURL] 按公开直链删除对象。
package provideraliyunoss
