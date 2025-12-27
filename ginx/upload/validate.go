package upload

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ValidateFormFile 校验表单文件（类型、大小），返回 FileHeader 和元信息。
//
// 如果字段为空（用户未上传），返回 nil, nil。
//
//	result, err := upload.ValidateFormFile(c, "avatar", 5*1024*1024, []string{"image/jpeg", "image/png"})
func ValidateFormFile(ctx *gin.Context, fieldName string, maxSize int64, allowedTypes []string) (*FormFileResult, error) {
	fh, err := ctx.FormFile(fieldName)
	if err != nil {
		// 字段为空或不存在，属于可选场景
		return nil, nil
	}

	if fh.Size > maxSize {
		return nil, fmt.Errorf("文件大小不能超过 %dMB", maxSize/(1024*1024))
	}

	contentType := fh.Header.Get("Content-Type")
	if !slices.Contains(allowedTypes, contentType) {
		return nil, errors.New("文件格式错误")
	}

	return &FormFileResult{
		FileHeader:  fh,
		Filename:    generateUniqueFilename(fh.Filename),
		ContentType: contentType,
	}, nil
}

// generateUniqueFilename 生成唯一文件名，格式: 时间戳_uuid.扩展名
func generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().Format("20060102150405")
	uniqueID := uuid.New().String()[:8]
	return timestamp + "_" + uniqueID + ext
}
