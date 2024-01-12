package middlewares

import (
	"gin-fleamarket/services"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService services.IAuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return
		}
		if !strings.HasPrefix(header, "Bearer ") {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		user, err := authService.GetUserFromToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}

		ctx.Set("user", user)

		ctx.Next()
	}
}