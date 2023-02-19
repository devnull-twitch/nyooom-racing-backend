package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type EditorLogin struct {
	Username string
	Password string
}

func GetEditorMiddleware(editors []*EditorLogin) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, password, ok := ctx.Request.BasicAuth()

		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		for _, editor := range editors {
			if editor.Username == username && editor.Password == password {
				return
			}
		}

		ctx.AbortWithStatus(http.StatusForbidden)
	}
}
