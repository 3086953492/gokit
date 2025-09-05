package logger

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// CleanupOldLogs 清理过期的日志文件
func CleanupOldLogs(filename string, maxAge, maxBackups int) error {
	dir := filepath.Dir(filename)
	ext := filepath.Ext(filename)
	base := filepath.Base(filename[:len(filename)-len(ext)])

	// 读取目录中的所有文件
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var logFiles []os.FileInfo

	// 筛选出相关的日志文件
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		// 匹配格式：basename-YYYY-MM-DD.ext 或 basename-YYYY-MM-DD.ext.gz
		if strings.HasPrefix(name, base+"-") && (strings.HasSuffix(name, ext) || strings.HasSuffix(name, ext+".gz")) {
			if info, err := file.Info(); err == nil {
				logFiles = append(logFiles, info)
			}
		}
	}

	// 按修改时间排序（最新的在前面）
	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].ModTime().After(logFiles[j].ModTime())
	})

	// 删除超过最大备份数量的文件
	if maxBackups > 0 && len(logFiles) > maxBackups {
		for i := maxBackups; i < len(logFiles); i++ {
			filePath := filepath.Join(dir, logFiles[i].Name())
			os.Remove(filePath)
		}
		logFiles = logFiles[:maxBackups]
	}

	// 删除超过最大年龄的文件
	if maxAge > 0 {
		cutoff := time.Now().AddDate(0, 0, -maxAge)
		for _, file := range logFiles {
			if file.ModTime().Before(cutoff) {
				filePath := filepath.Join(dir, file.Name())
				os.Remove(filePath)
			}
		}
	}

	return nil
}
