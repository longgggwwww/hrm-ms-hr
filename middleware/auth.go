package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims defines JWT claims with permissions
type CustomClaims struct {
	PermCodes []string `json:"perm_codes"`
	jwt.RegisteredClaims
}

// ValidateToken validates JWT token using JWT_SECRET from .env
func ValidateToken(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	_, _, err := new(jwt.Parser).ParseUnverified(tokenString, claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// AuthMiddleware checks JWT token and required permissions
func AuthMiddleware(requiredPerms []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "#1 AuthMiddleware: Unauthorized", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateToken(token)
		if err != nil {
			http.Error(w, "#2 AuthMiddleware: Invalid token", http.StatusUnauthorized)
			return
		}
		permMap := make(map[string]bool)
		for _, p := range claims.PermCodes {
			permMap[p] = true
		}
		for _, rp := range requiredPerms {
			fmt.Print(permMap[rp], 123, rp)
			if !permMap[rp] {
				http.Error(w, "#3 AuthMiddleware: Forbidden", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
