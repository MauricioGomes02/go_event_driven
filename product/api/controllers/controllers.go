package controllers

import (
	fmt "fmt"
	http "net/http"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	ProductController ProductController
}

func CreatedAtRoute(id any, output any, context *gin.Context) {
	path := context.Request.URL.Path
	context.Header("Location", fmt.Sprintf("%s/%s", path, id))
	context.JSON(http.StatusCreated, output)
}
