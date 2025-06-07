package ports

import "context"

type Logger interface {
	LogInformation(_context context.Context, message string, fields ...Field)
	LogWarning(_context context.Context, message string, fields ...Field)
	LogError(_context context.Context, message string, fields ...Field)
	With(_context context.Context, fields ...Field) context.Context
}

type Field struct {
	Key   string
	Value any
}
