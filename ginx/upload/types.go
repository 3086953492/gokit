package upload

import "mime/multipart"

// FormFileResult 保存校验后的文件元信息（不包含打开的句柄）。
type FormFileResult struct {
	FileHeader  *multipart.FileHeader
	Filename    string // 生成的唯一文件名
	ContentType string
}

