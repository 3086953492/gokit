package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// ObjectStorage 对象存储统一接口
type ObjectStorage interface {
	// Upload 上传文件并返回结果
	Upload(ctx context.Context, file FileObject) (UploadResult, error)
}

// KeyStrategy 对象 Key 生成策略接口
type KeyStrategy interface {
	Generate(file FileObject) string
}

// DatePathRandomKeyStrategy 默认策略：日期路径 + 随机串 + 扩展名
type DatePathRandomKeyStrategy struct{}

// Generate 生成对象 Key，格式：YYYY/MM/DD/YYYYMMDD_<随机串>.<扩展名>
func (s *DatePathRandomKeyStrategy) Generate(file FileObject) string {
	now := time.Now()
	datePath := now.Format("2006/01/02")
	datePrefix := now.Format("20060102")

	// 生成随机串
	randBytes := make([]byte, 8)
	rand.Read(randBytes)
	randStr := hex.EncodeToString(randBytes)

	// 推断扩展名
	ext := s.inferExtension(file)

	filename := fmt.Sprintf("%s_%s%s", datePrefix, randStr, ext)
	return filepath.ToSlash(filepath.Join(datePath, filename))
}

// inferExtension 推断扩展名
func (s *DatePathRandomKeyStrategy) inferExtension(file FileObject) string {
	// 优先从原始文件名推断
	if file.Filename != "" {
		ext := filepath.Ext(file.Filename)
		if ext != "" {
			return strings.ToLower(ext)
		}
	}

	// 从 ContentType 推断常见扩展名
	if file.ContentType != "" {
		ext := MimeToExtension(file.ContentType)
		if ext != "" {
			return ext
		}
	}

	return ""
}

// MimeToExtension MIME 类型到扩展名的映射（导出供子包复用）
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
