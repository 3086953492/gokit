package storage

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

// KeyGenerator 对象 Key 生成器接口。
type KeyGenerator interface {
	Generate(filename, contentType string) string
}

// DatePathKeyGenerator 日期路径 + 随机串 + 扩展名的 Key 生成器。
// 生成格式：YYYY/MM/DD/YYYYMMDD_<随机串>.<扩展名>
type DatePathKeyGenerator struct{}

// Generate 根据文件名和内容类型生成对象 Key。
func (g *DatePathKeyGenerator) Generate(filename, contentType string) string {
	now := time.Now()
	datePath := now.Format("2006/01/02")
	datePrefix := now.Format("20060102")

	randBytes := make([]byte, 8)
	rand.Read(randBytes)
	randStr := hex.EncodeToString(randBytes)

	ext := g.inferExtension(filename, contentType)
	name := fmt.Sprintf("%s_%s%s", datePrefix, randStr, ext)

	return filepath.ToSlash(filepath.Join(datePath, name))
}

// inferExtension 推断扩展名。
func (g *DatePathKeyGenerator) inferExtension(filename, contentType string) string {
	if filename != "" {
		ext := filepath.Ext(filename)
		if ext != "" {
			return strings.ToLower(ext)
		}
	}

	if contentType != "" {
		ext := MimeToExtension(contentType)
		if ext != "" {
			return ext
		}
	}

	return ""
}

var mimeToExt = map[string]string{
	"image/jpeg":       ".jpg",
	"image/png":        ".png",
	"image/gif":        ".gif",
	"image/webp":       ".webp",
	"image/svg+xml":    ".svg",
	"application/pdf":  ".pdf",
	"text/plain":       ".txt",
	"text/html":        ".html",
	"text/css":         ".css",
	"text/javascript":  ".js",
	"application/json": ".json",
	"application/xml":  ".xml",
	"application/zip":  ".zip",
	"video/mp4":        ".mp4",
	"audio/mpeg":       ".mp3",
}

// EscapeKey 对对象 key 进行 URL 路径编码，保留 "/" 作为目录分隔符。
func EscapeKey(key string) string {
	parts := strings.Split(key, "/")
	for i, p := range parts {
		parts[i] = url.PathEscape(p)
	}
	return strings.Join(parts, "/")
}

// UnescapeKey 对 URL 路径逐段解码，还原对象 key。
func UnescapeKey(escapedPath string) (string, error) {
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

// MimeToExtension 将 MIME 类型转换为文件扩展名。
// 支持带参数的 MIME 类型（如 "text/plain; charset=utf-8"），未知类型返回空串。
func MimeToExtension(mimeType string) string {
	baseMime, _, _ := strings.Cut(mimeType, ";")
	baseMime = strings.TrimSpace(baseMime)

	if ext, ok := mimeToExt[baseMime]; ok {
		return ext
	}
	return ""
}
