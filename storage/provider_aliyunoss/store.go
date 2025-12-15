package provider_aliyunoss

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"

	"github.com/3086953492/gokit/storage"
)

// ProviderName 阿里云 OSS 的 Provider 名称。
const ProviderName = "aliyun_oss"

// Store 阿里云 OSS 存储实现。
type Store struct {
	client   *oss.Client
	bucket   string
	endpoint string
	domain   string
}

// New 创建阿里云 OSS 存储实现。
func New(cfg Config) (storage.Store, error) {
	if cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" {
		return nil, fmt.Errorf("%w: AccessKeyID and AccessKeySecret are required", storage.ErrInvalidConfig)
	}
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("%w: Endpoint is required", storage.ErrInvalidConfig)
	}
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("%w: Bucket is required", storage.ErrInvalidConfig)
	}

	region := extractRegionFromEndpoint(cfg.Endpoint)
	credProvider := credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.AccessKeySecret)

	ossCfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credProvider).
		WithRegion(region).
		WithEndpoint(cfg.Endpoint)

	ossClient := oss.NewClient(ossCfg)

	return &Store{
		client:   ossClient,
		bucket:   cfg.Bucket,
		endpoint: cfg.Endpoint,
		domain:   cfg.Domain,
	}, nil
}

// Upload 上传对象到阿里云 OSS。
func (s *Store) Upload(ctx context.Context, key string, r io.Reader, opts *storage.WriteOptions) (*storage.ObjectMeta, error) {
	putReq := &oss.PutObjectRequest{
		Bucket: oss.Ptr(s.bucket),
		Key:    oss.Ptr(key),
		Body:   r,
	}

	if opts != nil {
		if opts.ContentType != "" {
			putReq.ContentType = oss.Ptr(opts.ContentType)
		}
		if opts.CacheControl != "" {
			putReq.CacheControl = oss.Ptr(opts.CacheControl)
		}
		if opts.ContentLength > 0 {
			putReq.ContentLength = oss.Ptr(opts.ContentLength)
		}
		if len(opts.UserMeta) > 0 {
			putReq.Metadata = opts.UserMeta
		}
	}

	result, err := s.client.PutObject(ctx, putReq)
	if err != nil {
		return nil, fmt.Errorf("aliyun oss upload %s: %w", key, err)
	}

	meta := &storage.ObjectMeta{
		Key:         key,
		ContentType: safeDeref(opts, func(o *storage.WriteOptions) string { return o.ContentType }),
		URL:         s.objectURL(key),
	}
	if result.ETag != nil {
		meta.ETag = *result.ETag
	}

	return meta, nil
}

// Download 从阿里云 OSS 下载对象。
func (s *Store) Download(ctx context.Context, key string, opts *storage.ReadOptions) (io.ReadCloser, *storage.ObjectMeta, error) {
	getReq := &oss.GetObjectRequest{
		Bucket: oss.Ptr(s.bucket),
		Key:    oss.Ptr(key),
	}

	if opts != nil && opts.Range != "" {
		getReq.Range = oss.Ptr(opts.Range)
	}

	result, err := s.client.GetObject(ctx, getReq)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil, fmt.Errorf("aliyun oss download %s: %w", key, storage.ErrNotFound)
		}
		return nil, nil, fmt.Errorf("aliyun oss download %s: %w", key, err)
	}

	meta := &storage.ObjectMeta{
		Key:         key,
		Size:        safeDerefInt64(result.ContentLength),
		ContentType: safeDerefString(result.ContentType),
		ETag:        safeDerefString(result.ETag),
		URL:         s.objectURL(key),
	}
	if result.LastModified != nil {
		meta.LastModified = *result.LastModified
	}

	return result.Body, meta, nil
}

// Delete 删除阿里云 OSS 中的对象。
func (s *Store) Delete(ctx context.Context, key string, opts *storage.DeleteOptions) error {
	delReq := &oss.DeleteObjectRequest{
		Bucket: oss.Ptr(s.bucket),
		Key:    oss.Ptr(key),
	}

	_, err := s.client.DeleteObject(ctx, delReq)
	if err != nil {
		if isNotFoundError(err) {
			return fmt.Errorf("aliyun oss delete %s: %w", key, storage.ErrNotFound)
		}
		return fmt.Errorf("aliyun oss delete %s: %w", key, err)
	}

	return nil
}

// List 列举阿里云 OSS 中指定前缀的对象。
func (s *Store) List(ctx context.Context, prefix string, opts *storage.ListOptions) (*storage.ListResult, error) {
	listReq := &oss.ListObjectsV2Request{
		Bucket: oss.Ptr(s.bucket),
		Prefix: oss.Ptr(prefix),
	}

	if opts != nil {
		if opts.MaxKeys > 0 {
			listReq.MaxKeys = int32(opts.MaxKeys)
		}
		if opts.Marker != "" {
			listReq.ContinuationToken = oss.Ptr(opts.Marker)
		}
		if opts.Delimiter != "" {
			listReq.Delimiter = oss.Ptr(opts.Delimiter)
		}
	}

	result, err := s.client.ListObjectsV2(ctx, listReq)
	if err != nil {
		return nil, fmt.Errorf("aliyun oss list prefix=%s: %w", prefix, err)
	}

	objects := make([]*storage.ObjectMeta, 0, len(result.Contents))
	for _, obj := range result.Contents {
		objKey := safeDerefString(obj.Key)
		meta := &storage.ObjectMeta{
			Key:  objKey,
			Size: safeDerefInt64(obj.Size),
			ETag: safeDerefString(obj.ETag),
			URL:  s.objectURL(objKey),
		}
		if obj.LastModified != nil {
			meta.LastModified = *obj.LastModified
		}
		objects = append(objects, meta)
	}

	commonPrefixes := make([]string, 0, len(result.CommonPrefixes))
	for _, cp := range result.CommonPrefixes {
		if cp.Prefix != nil {
			commonPrefixes = append(commonPrefixes, *cp.Prefix)
		}
	}

	listResult := &storage.ListResult{
		Objects:        objects,
		IsTruncated:    result.IsTruncated,
		CommonPrefixes: commonPrefixes,
	}
	if result.NextContinuationToken != nil {
		listResult.NextMarker = *result.NextContinuationToken
	}

	return listResult, nil
}

// Exists 检查阿里云 OSS 中对象是否存在。
func (s *Store) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := s.client.IsObjectExist(ctx, s.bucket, key)
	if err != nil {
		return false, fmt.Errorf("aliyun oss exists %s: %w", key, err)
	}
	return exists, nil
}

// Head 获取阿里云 OSS 中对象的元信息。
func (s *Store) Head(ctx context.Context, key string) (*storage.ObjectMeta, error) {
	headReq := &oss.HeadObjectRequest{
		Bucket: oss.Ptr(s.bucket),
		Key:    oss.Ptr(key),
	}

	result, err := s.client.HeadObject(ctx, headReq)
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("aliyun oss head %s: %w", key, storage.ErrNotFound)
		}
		return nil, fmt.Errorf("aliyun oss head %s: %w", key, err)
	}

	meta := &storage.ObjectMeta{
		Key:         key,
		Size:        safeDerefInt64(result.ContentLength),
		ContentType: safeDerefString(result.ContentType),
		ETag:        safeDerefString(result.ETag),
		URL:         s.objectURL(key),
	}
	if result.LastModified != nil {
		meta.LastModified = *result.LastModified
	}

	return meta, nil
}

// extractRegionFromEndpoint 从 endpoint 提取 region。
func extractRegionFromEndpoint(endpoint string) string {
	endpoint = strings.TrimPrefix(endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")

	if strings.HasPrefix(endpoint, "oss-") {
		parts := strings.Split(endpoint, ".")
		if len(parts) > 0 {
			return strings.TrimPrefix(parts[0], "oss-")
		}
	}

	return endpoint
}

// isNotFoundError 判断错误是否为对象不存在。
func isNotFoundError(err error) bool {
	var serviceErr *oss.ServiceError
	if errors.As(err, &serviceErr) {
		return serviceErr.Code == "NoSuchKey" || serviceErr.StatusCode == 404
	}
	return false
}

// safeDeref 安全解引用 WriteOptions 并获取字段值。
func safeDeref(opts *storage.WriteOptions, getter func(*storage.WriteOptions) string) string {
	if opts == nil {
		return ""
	}
	return getter(opts)
}

// safeDerefString 安全解引用 string 指针。
func safeDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// safeDerefInt64 返回安全的 int64 值。
// 当前 OSS SDK 中相关字段为非指针类型，这里预留一层封装便于未来扩展。
func safeDerefInt64(v int64) int64 {
	return v
}

// objectURL 生成对象的公开可访问 URL。
// 优先使用 domain（配置的自定义域名或 CDN 域名），否则使用 bucket.endpoint 拼接。
func (s *Store) objectURL(key string) string {
	baseURL := s.baseURL()
	escapedKey := escapeKey(key)
	return baseURL + "/" + escapedKey
}

// baseURL 返回拼接 URL 的基础域名（带 scheme，无尾部 /）。
func (s *Store) baseURL() string {
	if s.domain != "" {
		return normalizeBaseURL(s.domain)
	}
	endpoint := normalizeEndpoint(s.endpoint)
	return "https://" + s.bucket + "." + endpoint
}

// normalizeBaseURL 规范化自定义域名，确保带 https:// 且无尾部 /。
func normalizeBaseURL(domain string) string {
	d := strings.TrimSpace(domain)
	if !strings.HasPrefix(d, "http://") && !strings.HasPrefix(d, "https://") {
		d = "https://" + d
	}
	return strings.TrimSuffix(d, "/")
}

// normalizeEndpoint 去除 endpoint 的 scheme 前缀和尾部 /。
func normalizeEndpoint(endpoint string) string {
	e := strings.TrimSpace(endpoint)
	e = strings.TrimPrefix(e, "https://")
	e = strings.TrimPrefix(e, "http://")
	return strings.TrimSuffix(e, "/")
}

// escapeKey 对 key 进行 URL 编码，保留 "/" 作为目录分隔符。
func escapeKey(key string) string {
	parts := strings.Split(key, "/")
	for i, p := range parts {
		parts[i] = url.PathEscape(p)
	}
	return strings.Join(parts, "/")
}

// ---------------------------------------------------------------------------
// URLKeyResolver 接口实现（可选能力，用于 Manager.DeleteByURL）
// ---------------------------------------------------------------------------

// AllowedHosts 返回当前 Store 允许的域名列表（仅 host 部分，不含 scheme）。
// 包含自定义域名（如配置）以及默认的 bucket.endpoint host。
func (s *Store) AllowedHosts() []string {
	hosts := make([]string, 0, 2)

	// 默认 host：bucket.endpoint（endpoint 已规范化去除 scheme）
	defaultHost := s.bucket + "." + normalizeEndpoint(s.endpoint)
	hosts = append(hosts, defaultHost)

	// 自定义域名（如配置）
	if s.domain != "" {
		customHost := extractHost(s.domain)
		if customHost != "" && customHost != defaultHost {
			hosts = append(hosts, customHost)
		}
	}

	return hosts
}

// KeyFromURL 从已解析的 URL 提取对象 key。
// 仅支持 objectURL() 生成的格式（/{escapedKey}）；返回的 key 已做 URL 解码。
func (s *Store) KeyFromURL(u *url.URL) (string, error) {
	path := u.Path
	if path == "" || path == "/" {
		return "", fmt.Errorf("%w: empty path", storage.ErrInvalidURL)
	}

	// 去除前导 /
	if path[0] == '/' {
		path = path[1:]
	}

	// 逐段解码
	key, err := unescapeKey(path)
	if err != nil {
		return "", fmt.Errorf("%w: %v", storage.ErrInvalidURL, err)
	}

	if key == "" {
		return "", fmt.Errorf("%w: empty key after decode", storage.ErrInvalidURL)
	}

	return key, nil
}

// extractHost 从带或不带 scheme 的域名字符串中提取 host 部分（不含端口）。
func extractHost(domain string) string {
	d := strings.TrimSpace(domain)
	d = strings.TrimPrefix(d, "https://")
	d = strings.TrimPrefix(d, "http://")
	d = strings.TrimSuffix(d, "/")

	// 去除端口（如有）
	if idx := strings.LastIndex(d, ":"); idx != -1 {
		// 确保不是 IPv6 地址中的冒号
		if !strings.Contains(d, "[") {
			d = d[:idx]
		}
	}
	return d
}

// unescapeKey 对 URL 路径逐段解码，还原对象 key。
func unescapeKey(escapedPath string) (string, error) {
	parts := strings.Split(escapedPath, "/")
	for i, p := range parts {
		decoded, err := url.PathUnescape(p)
		if err != nil {
			return "", err
		}
		parts[i] = decoded
	}
	return strings.Join(parts, "/"), nil
}
