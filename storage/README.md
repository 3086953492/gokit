# storage

对象存储抽象层：用统一的 `Store` 接口与 `Manager` 入口封装上传、下载、列举、删除等操作；当前内置阿里云 OSS 与本地文件系统实现。

## 安装

```bash
go get github.com/3086953492/gokit/storage
```

使用具体后端时需同时引入对应子包：

```bash
go get github.com/3086953492/gokit/storage/provider_aliyunoss
go get github.com/3086953492/gokit/storage/provider_local
```

## 核心概念

| 类型 | 说明 |
|------|------|
| `Store` | 后端接口：`Upload` / `Download` / `Delete` / `List` / `Exists` / `Head` |
| `Manager` | 对外统一入口，线程安全；校验 key、包装错误、保证 `List` 返回非 nil 切片 |
| `ObjectMeta` | 对象元信息（含可选公开直链 `URL`） |
| `URLKeyResolver` | 可选接口：实现后 `Manager.DeleteByURL` 可根据直链解析并删除对象 |

## 快速开始（阿里云 OSS）

```go
import (
	"context"
	"strings"

	"github.com/3086953492/gokit/storage"
	"github.com/3086953492/gokit/storage/provider_aliyunoss"
)

func main() {
	store, err := provider_aliyunoss.New(provider_aliyunoss.Config{
		AccessKeyID:     "...",
		AccessKeySecret: "...",
		Endpoint:        "oss-cn-hangzhou.aliyuncs.com",
		Bucket:          "your-bucket",
		Domain:          "", // 可选：自定义域名或 CDN，用于 ObjectMeta.URL
	})
	if err != nil {
		panic(err)
	}

	mgr, err := storage.NewManager(storage.WithStore(store))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	key := "uploads/demo.txt"

	meta, err := mgr.Upload(ctx, key, strings.NewReader("hello"),
		storage.WithContentType("text/plain; charset=utf-8"),
	)
	// meta.URL 为公开访问 URL（需 Bucket/权限允许匿名读）

	rc, dlMeta, err := mgr.Download(ctx, key)
	if err != nil {
		panic(err)
	}
	defer rc.Close()
	_ = dlMeta
}
```

## 快速开始（本地文件系统）

```go
import (
	"context"
	"strings"

	"github.com/3086953492/gokit/storage"
	providerlocal "github.com/3086953492/gokit/storage/provider_local"
)

func main() {
	store, err := providerlocal.New(providerlocal.Config{
		Root:    "./data/storage",
		BaseURL: "https://static.example.com/files", // 可选：为空时不生成 ObjectMeta.URL
	})
	if err != nil {
		panic(err)
	}

	mgr, err := storage.NewManager(storage.WithStore(store))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	meta, err := mgr.Upload(ctx, "avatars/user-1.txt", strings.NewReader("hello"),
		storage.WithContentType("text/plain; charset=utf-8"),
	)
	if err != nil {
		panic(err)
	}

	_ = meta.URL // 配置 BaseURL 时才会生成公开直链
}
```

## Manager 与选项

- **创建**：`NewManager` 必须配合 `WithStore(store)`，否则返回 `ErrInvalidConfig`。
- **上传**：`WithContentType`、`WithCacheControl`、`WithContentLength`（部分后端需要已知长度）、`WithUserMeta`。
- **下载**：`WithRange`（如 `bytes=0-1023`）。
- **列举**：默认单次最多 1000 条；`WithMaxKeys`、`WithMarker`（续页）、`WithDelimiter`（如 `/` 模拟目录）。
- **底层实现**：`Manager.Store()` 一般无需使用。

## 本地存储说明

- **根目录**：`provider_local.Config.Root` 为必填项，所有对象都会落在该目录下。
- **Key 规则**：逻辑 key 使用 `/` 作为分隔符；本地实现会拒绝空 key、绝对路径、`..` 路径穿越和 `\\` 分隔符。
- **URL 能力**：仅在配置 `BaseURL` 时生成 `ObjectMeta.URL`，并支持 `DeleteByURL`；未配置时会返回 `ErrURLDeleteUnsupported`。
- **范围下载**：本地实现支持单段 `bytes=` 范围读取。
- **权限与元数据**：可通过 `DirPerm`、`FilePerm` 控制自动创建目录和写入文件权限；`WriteOptions.UserMeta`、`CacheControl` 不会持久化到本地文件系统。
- **列举行为**：`List` 支持 `Prefix`、`MaxKeys`、`Marker`、`Delimiter`；分页基于稳定排序后的 key / 公共前缀标记。

## 按 URL 删除

若 `Store` 实现了 `URLKeyResolver`（阿里云 OSS 与已配置 `BaseURL` 的本地存储均支持），可对 `ObjectMeta.URL` 形式的公开直链调用：

```go
err := mgr.DeleteByURL(ctx, meta.URL)
```

约束：仅支持 `http`/`https`；域名须在 `AllowedHosts()` 内；路径须为本库生成 URL 的格式。未实现接口时返回 `ErrURLDeleteUnsupported`。

## Key 与工具

- **Key 校验**：`Manager` 对 key 非空校验，空 key 返回 `ErrInvalidKey`。
- **`KeyGenerator`**：可自定义对象 key 规则。
- **`DatePathKeyGenerator`**：生成 `YYYY/MM/DD/YYYYMMDD_<随机16hex>.<扩展名>`，扩展名来自文件名或 `MimeToExtension`。
- **`MimeToExtension`**：常见 MIME 到扩展名映射（含带 `; charset=` 的 MIME）。

## 错误处理

包内预定义错误可通过 `errors.Is` 判断，例如：

- `ErrNotFound`、`ErrAlreadyExists`
- `ErrInvalidKey`、`ErrInvalidConfig`、`ErrInvalidURL`
- `ErrDomainNotAllowed`、`ErrURLDeleteUnsupported`
- `ErrBackendUnavailable`、`ErrPermissionDenied`

`Manager` 方法在包装底层错误时通常使用 `%w`，可与 `errors.Is` / `errors.As` 链式使用。

## 扩展其他后端

实现 `storage.Store` 即可接入 `NewManager(storage.WithStore(...))`。若需 `DeleteByURL`，同时实现 `storage.URLKeyResolver`。

## 依赖说明

`provider_aliyunoss` 基于 [alibabacloud-oss-go-sdk-v2](https://github.com/aliyun/alibabacloud-oss-go-sdk-v2)。
