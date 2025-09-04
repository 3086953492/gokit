package configs

import "time"

type OAuthConfig struct {
	// 授权服务器基本配置
	Issuer      string `mapstructure:"issuer"`
	AuthorizeUI string `mapstructure:"authorize_ui"` // 授权页面URL

	// 授权码配置
	AuthorizationCodeTTL time.Duration `mapstructure:"authorization_code_ttl"`

	// 访问令牌配置
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`

	// 支持的授权类型
	SupportedGrantTypes []string `mapstructure:"supported_grant_types"`

	// 支持的响应类型
	SupportedResponseTypes []string `mapstructure:"supported_response_types"`

	// 默认权限范围
	DefaultScopes []string `mapstructure:"default_scopes"`

	// PKCE 支持
	RequirePKCE              bool     `mapstructure:"require_pkce"`
	SupportedChallengeMethod []string `mapstructure:"supported_challenge_methods"`
}
