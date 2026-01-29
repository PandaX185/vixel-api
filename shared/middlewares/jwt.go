package middlewares

import (
	"strings"
	"vixel/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func validateAccessToken(header string) (*jwt.Token, error) {
	token := strings.TrimPrefix(header, "Bearer ")
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.JWTSecret), nil
	})
}

func JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Authorization header missing"})
			return
		}

		token, err := validateAccessToken(authHeader)
		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Invalid token claims"})
			return
		}

		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Invalid user ID in token"})
			return
		}

		userID := uint(userIDFloat)
		ctx.Set("user_id", userID)
		ctx.Next()
	}
}
