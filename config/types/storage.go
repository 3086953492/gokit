package types

import (
	"fmt"
	"strings"
)

// 与 storage 子包 Provider 名称一致。
const (
	StorageProviderAliyunOSS = "aliyun_oss"
	StorageProviderLocal     = "local"
)

// LocalStorageConfig 本地文件系统存储子配置（对应 providerlocal.Config）。
// DirPerm、FilePerm 为 Unix 权限位十进制值（如 0o755=493）；零值表示使用 provider 默认值。
type LocalStorageConfig struct {
	Root     string `json:"root" yaml:"root" mapstructure:"root"`
	BaseURL  string `json:"base_url" yaml:"base_url" mapstructure:"base_url"`
	DirPerm  uint32 `json:"dir_perm,omitempty" yaml:"dir_perm,omitempty" mapstructure:"dir_perm"`
	FilePerm uint32 `json:"file_perm,omitempty" yaml:"file_perm,omitempty" mapstructure:"file_perm"`
}

// AliyunOSSConfig 阿里云 OSS 存储子配置（对应 provideraliyunoss.Config）。
type AliyunOSSConfig struct {
	AccessKeyID     string `json:"access_key_id" yaml:"access_key_id" mapstructure:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret" yaml:"access_key_secret" mapstructure:"access_key_secret"`
	Endpoint        string `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`
	Bucket          string `json:"bucket" yaml:"bucket" mapstructure:"bucket"`
	Domain          string `json:"domain" yaml:"domain" mapstructure:"domain"` // 可选：自定义域名
}


// StorageConfig 统一存储配置：由 Provider 选择后端，再读取对应嵌套子块。
type StorageConfig struct {
	Provider  string             `json:"provider" yaml:"provider" mapstructure:"provider"`
	AliyunOSS AliyunOSSConfig    `json:"aliyun_oss" yaml:"aliyun_oss" mapstructure:"aliyun_oss"`
	Local     LocalStorageConfig `json:"local" yaml:"local" mapstructure:"local"`
}

// Validate 校验 storage 段；Provider 为空视为未启用存储。
func (c *StorageConfig) Validate() error {
	switch strings.TrimSpace(c.Provider) {
	case "":
		return nil
	case StorageProviderAliyunOSS:
		if c.AliyunOSS.AccessKeyID == "" || c.AliyunOSS.AccessKeySecret == "" {
			return fmt.Errorf("storage: aliyun_oss requires access_key_id and access_key_secret")
		}
		if c.AliyunOSS.Endpoint == "" {
			return fmt.Errorf("storage: aliyun_oss requires endpoint")
		}
		if c.AliyunOSS.Bucket == "" {
			return fmt.Errorf("storage: aliyun_oss requires bucket")
		}
		return nil
	case StorageProviderLocal:
		if strings.TrimSpace(c.Local.Root) == "" {
			return fmt.Errorf("storage: local requires root")
		}
		return nil
	default:
		return fmt.Errorf("storage: unknown provider %q", c.Provider)
	}
}
