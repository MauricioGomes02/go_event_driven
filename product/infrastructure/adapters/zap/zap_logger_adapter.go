package zap

import (
	"context"
	"go_event_driven/product/domain/ports"
	"os"

	zapLogger "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	loggerFieldsContextKey = &struct{}{}
)

type ZapLogger struct {
	logger *zapLogger.Logger
}

func NewZapLogger() *ZapLogger {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:       "@timestamp",
		LevelKey:      "log.level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)

	logger := zapLogger.New(core)
	return &ZapLogger{logger: logger}
}

func (logger *ZapLogger) LogInformation(_context context.Context, message string, fields ...ports.Field) {
	logger.logWithContext(_context, "information", message, fields...)
}

func (logger *ZapLogger) LogWarning(_context context.Context, message string, fields ...ports.Field) {
	logger.logWithContext(_context, "warning", message, fields...)
}

func (logger *ZapLogger) LogError(_context context.Context, message string, fields ...ports.Field) {
	logger.logWithContext(_context, "error", message, fields...)
}

func (logger *ZapLogger) With(_context context.Context, fields ...ports.Field) context.Context {
	existingFields := []ports.Field{}
	existingValue := _context.Value(loggerFieldsContextKey)
	if existingValue != nil {
		existingFields = existingValue.([]ports.Field)
	}
	return context.WithValue(_context, loggerFieldsContextKey, append(existingFields, fields...))
}

func convertToZapFields(fields []ports.Field) []zapLogger.Field {
	_fields := make([]zapLogger.Field, 0, len(fields))
	for _, field := range fields {
		_fields = append(_fields, zapLogger.Any(field.Key, field.Value))
	}
	return _fields
}

func (logger *ZapLogger) logWithContext(_context context.Context, level string, message string, fields ...ports.Field) {
	contextFieldsValue := _context.Value(loggerFieldsContextKey)
	allFields := []zapLogger.Field{}
	if contextFieldsValue != nil {
		allFields = append(allFields, convertToZapFields(contextFieldsValue.([]ports.Field))...)
	}

	allFields = append(allFields, convertToZapFields(fields)...)

	switch level {
	case "information":
		logger.logger.Info(message, allFields...)
	case "warning":
		logger.logger.Warn(message, allFields...)
	case "error":
		logger.logger.Error(message, allFields...)
	}
}
