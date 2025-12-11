package aliyunoss

import (
	"context"
	"fmt"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"

	"github.com/3086953492/gokit/storage"
)

// ProviderName 阿里云 OSS 的 Provider 名称
const ProviderName = "aliyun_oss"

// Config 阿里云 OSS 配置
type Config struct {
	AccessKeyID     string // 阿里云 AccessKeyId
	AccessKeySecret string // 阿里云 AccessKeySecret
	Endpoint        string // 例如 oss-cn-hangzhou.aliyuncs.com
	Bucket          string // 默认使用的 Bucket 名称
	Domain          string // 自定义域名（可选），为空时使用默认规则拼接 URL
}

// Option 可选配置函数
type Option func(*Client)

// WithKeyStrategy 设置自定义的 Key 生成策略
func WithKeyStrategy(strategy storage.KeyStrategy) Option {
	return func(c *Client) {
		c.keyStrategy = strategy
	}
}

// Client 阿里云 OSS 客户端实现
type Client struct {
	client      *oss.Client
	bucket      string
	endpoint    string
	domain      string
	keyStrategy storage.KeyStrategy
}

// New 创建阿里云 OSS 客户端，返回 storage.ObjectStorage 接口
func New(cfg Config, opts ...Option) (storage.ObjectStorage, error) {
	// 参数校验
	if cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" {
		return nil, fmt.Errorf("aliyun oss: AccessKeyID and AccessKeySecret are required")
	}
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("aliyun oss: Endpoint is required")
	}
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("aliyun oss: Bucket is required")
	}

	// 从 endpoint 推断 region
	region := extractRegionFromEndpoint(cfg.Endpoint)

	// 创建凭证提供者
	credProvider := credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.AccessKeySecret)

	// 创建 OSS 配置
	ossCfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credProvider).
		WithRegion(region).
		WithEndpoint(cfg.Endpoint)

	// 创建 OSS 客户端
	ossClient := oss.NewClient(ossCfg)

	// 检查 Bucket 是否存在，不存在则创建
	ctx := context.Background()
	exists, err := ossClient.IsBucketExist(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("aliyun oss: failed to check bucket %s: %w", cfg.Bucket, err)
	}
	if !exists {
		// Bucket 不存在，尝试创建
		_, createErr := ossClient.PutBucket(ctx, &oss.PutBucketRequest{
			Bucket: oss.Ptr(cfg.Bucket),
		})
		if createErr != nil {
			return nil, fmt.Errorf("aliyun oss: failed to create bucket %s: %w", cfg.Bucket, createErr)
		}
	}

	c := &Client{
		client:      ossClient,
		bucket:      cfg.Bucket,
		endpoint:    cfg.Endpoint,
		domain:      cfg.Domain,
		keyStrategy: &storage.DatePathRandomKeyStrategy{}, // 使用核心包的默认策略
	}

	// 应用可选配置
	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// extractRegionFromEndpoint 从 endpoint 提取 region
// 例如 oss-cn-hangzhou.aliyuncs.com -> cn-hangzhou
func extractRegionFromEndpoint(endpoint string) string {
	// 去除协议前缀
	endpoint = strings.TrimPrefix(endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")

	// 提取 oss-<region>.aliyuncs.com 中的 region
	if strings.HasPrefix(endpoint, "oss-") {
		parts := strings.Split(endpoint, ".")
		if len(parts) > 0 {
			return strings.TrimPrefix(parts[0], "oss-")
		}
	}

	// 兜底返回 endpoint 本身
	return endpoint
}

// Upload 上传文件并返回结果
func (c *Client) Upload(ctx context.Context, file storage.FileObject) (storage.UploadResult, error) {
	// 生成对象 Key
	objectKey := c.keyStrategy.Generate(file)

	// 构建上传请求
	putReq := &oss.PutObjectRequest{
		Bucket:        oss.Ptr(c.bucket),
		Key:           oss.Ptr(objectKey),
		Body:          file.Reader,
		ContentLength: oss.Ptr(file.Size),
	}

	// 设置 Content-Type
	if file.ContentType != "" {
		putReq.ContentType = oss.Ptr(file.ContentType)
	}

	// 执行上传
	_, err := c.client.PutObject(ctx, putReq)
	if err != nil {
		return storage.UploadResult{}, fmt.Errorf("aliyun oss: failed to upload object %s: %w", objectKey, err)
	}

	// 拼接访问 URL
	url := c.buildURL(objectKey)

	return storage.UploadResult{
		Provider: ProviderName,
		Bucket:   c.bucket,
		Key:      objectKey,
		URL:      url,
	}, nil
}

// buildURL 拼接对象访问 URL
func (c *Client) buildURL(objectKey string) string {
	if c.domain != "" {
		// 使用自定义域名
		domain := strings.TrimSuffix(c.domain, "/")
		return fmt.Sprintf("%s/%s", domain, objectKey)
	}

	// 使用默认 OSS 域名
	endpoint := strings.TrimPrefix(c.endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")
	return fmt.Sprintf("https://%s.%s/%s", c.bucket, endpoint, objectKey)
}
