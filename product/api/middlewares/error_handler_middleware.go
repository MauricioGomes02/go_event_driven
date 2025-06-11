package middlewares

import (
	"context"
	outputmodels "go_event_driven/product/api/output_models"
	"go_event_driven/product/domain/ports"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleErrorMiddleware(logger ports.Logger) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		defer func() {
			_recover := recover()
			if _recover != nil {
				getValue, _ := ginContext.Get("contextLogger")
				contextLogger := getValue.(context.Context)

				_error, ok := _recover.(error)
				detail := ""
				if ok {
					detail = _error.Error()
				}
				problemDetail := &outputmodels.ProblemDetail{
					Title:    "An exception was thrown",
					Detail:   detail,
					Status:   http.StatusInternalServerError,
					Instance: ginContext.Request.URL.Path,
				}

				logger.LogError(
					contextLogger,
					"Error processing request",
					ports.Field{Key: "error.reason", Value: detail})

				ginContext.AbortWithStatusJSON(http.StatusInternalServerError, problemDetail)
			}
		}()

		ginContext.Next()
	}
}

func HandleLoggingMiddleware(_context context.Context, logger ports.Logger) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		start := time.Now()
		_context = logger.With(
			_context,
			ports.Field{Key: "http.method", Value: ginContext.Request.Method},
			ports.Field{Key: "http.path", Value: ginContext.Request.URL.Path},
			ports.Field{Key: "http.client.ip", Value: ginContext.ClientIP()},
			ports.Field{Key: "http.user.agent", Value: ginContext.Request.UserAgent()},
		)

		ginContext.Set("contextLogger", _context)
		ginContext.Next()
		statusCode := ginContext.Writer.Status()
		latency := time.Since(start)

		logger.LogInformation(
			_context,
			"Request completed",
			ports.Field{Key: "http.status", Value: statusCode},
			ports.Field{Key: "http.duration", Value: latency})
	}
}
