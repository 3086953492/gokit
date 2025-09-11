package oauth

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateAuthorizationCode 生成授权码
func GenerateAuthorizationCode() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
