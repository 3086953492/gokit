package logger

import "errors"

var (
	// ErrNoOutput 未配置任何输出
	ErrNoOutput = errors.New("logger: no output configured (enable console or file)")
	// ErrEmptyFilename 文件输出未指定文件名
	ErrEmptyFilename = errors.New("logger: file output enabled but filename is empty")
)

