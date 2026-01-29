package user

import (
	"time"
	"vixel/config"

	"github.com/golang-jwt/jwt/v5"
)

func generateAccessToken(user *User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   jwt.NewNumericDate(time.Now().Local().Add(time.Hour)),
	})

	return token.SignedString([]byte(config.Config.JWTSecret))
}

func extractUserIDFromToken(tokenStr string) (uint, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.JWTSecret), nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(float64); ok {
			return uint(sub), nil
		}
	}

	return 0, nil
}
