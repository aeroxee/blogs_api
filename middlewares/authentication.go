package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/aZ4ziL/blogs_api/auth"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authenticationHeader := ctx.Request.Header.Get("Authorization")
		if !strings.Contains(authenticationHeader, "Bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "autentikasi di butuhkan.",
			})
			return
		}
		token := strings.Replace(authenticationHeader, "Bearer ", "", -1)
		claims, err := auth.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		newContext := context.WithValue(context.Background(), &auth.UserAuth{}, claims)
		ctx.Request = ctx.Request.WithContext(newContext)
		ctx.Next()
	}
}
