package storage

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
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

// MimeToExtension 将 MIME 类型转换为文件扩展名。
func MimeToExtension(mimeType string) string {
	mimeMap := map[string]string{
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

	// 处理带参数的 MIME 类型，例如 "text/plain; charset=utf-8"
	baseMime := strings.Split(mimeType, ";")[0]
	baseMime = strings.TrimSpace(baseMime)

	if ext, ok := mimeMap[baseMime]; ok {
		return ext
	}
	return ""
}
