package middlewares

import (
	outputmodels "go_event_driven/product/api/output_models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleErrorMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			_recover := recover()
			if _recover != nil {
				_error, ok := _recover.(error)
				detail := ""
				if ok {
					detail = _error.Error()
				}
				problemDetail := &outputmodels.ProblemDetail{
					Title:    "An exception was thrown",
					Detail:   detail,
					Status:   http.StatusInternalServerError,
					Instance: context.Request.URL.Path,
				}
				context.AbortWithStatusJSON(http.StatusInternalServerError, problemDetail)
			}
		}()

		context.Next()
	}
}
