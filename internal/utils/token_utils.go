package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ExtractUserIDFromToken extracts user_id from JWT token in Authorization header
func ExtractUserIDFromToken(c *gin.Context) (int, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("authorization header missing")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return 0, fmt.Errorf("invalid authorization header format")
	}
	tokenString := parts[1]

	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	log.Println("claims", claims["user_id"])

	var userID int
	switch v := claims["user_id"].(type) {
	case string:
		var err error
		userID, err = strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("invalid user_id in token")
		}
	case float64:
		userID = int(v)
	case int:
		userID = v
	default:
		return 0, fmt.Errorf("user_id not found in token or has invalid type")
	}

	return userID, nil
}
