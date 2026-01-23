package i18n

import "errors"

var (
	// ErrMessageNotFound 翻译消息未找到
	ErrMessageNotFound = errors.New("message not found")
	// ErrInvalidConfig 配置无效
	ErrInvalidConfig = errors.New("invalid config")
	// ErrLoadFailed 加载翻译文件失败
	ErrLoadFailed = errors.New("load message file failed")
)
