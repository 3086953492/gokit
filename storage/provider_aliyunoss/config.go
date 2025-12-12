package provider_aliyunoss

// Config 阿里云 OSS 配置。
type Config struct {
	AccessKeyID     string // 阿里云 AccessKeyId
	AccessKeySecret string // 阿里云 AccessKeySecret
	Endpoint        string // 例如 oss-cn-hangzhou.aliyuncs.com
	Bucket          string // 默认使用的 Bucket 名称
	Domain          string // 自定义域名（可选），为空时使用默认规则拼接 URL
}
