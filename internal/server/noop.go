package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetNoopHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	}
}
