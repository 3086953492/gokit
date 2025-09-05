package logger

import (
	"time"

	"go.uber.org/zap"
)

func LogError(funcName, operation, message string, err error, fields ...zap.Field) {
	baseFields := []zap.Field{
		zap.String("function", funcName),
		zap.String("operation", operation),
		zap.Time("timestamp", time.Now()),
		zap.Error(err),
	}

	allFields := append(baseFields, fields...)
	Error(message, allFields...)
}