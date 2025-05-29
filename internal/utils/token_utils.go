package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ExtractIDsFromToken extracts user_id, org_id, employee_id from JWT token in Authorization header
func ExtractIDsFromToken(c *gin.Context) (map[string]int, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header missing")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, fmt.Errorf("invalid authorization header format")
	}
	tokenString := parts[1]

	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	result := make(map[string]int)
	for _, key := range []string{"user_id", "org_id", "employee_id"} {
		if val, ok := claims[key]; ok {
			switch v := val.(type) {
			case string:
				num, err := strconv.Atoi(v)
				if err != nil {
					return nil, fmt.Errorf("invalid %s in token", key)
				}
				result[key] = num
			case float64:
				result[key] = int(v)
			case int:
				result[key] = v
			}
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no valid id found in token")
	}
	return result, nil
}
