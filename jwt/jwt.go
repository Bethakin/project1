package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing Authorization header"})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := ValidateJWT(tokenString, secret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			sub, ok := claims["sub"].(float64)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token payload"})
			}
			userIDFromToken := fmt.Sprintf("%.0f", sub)
			userIDParam := c.Param("users_id")
			if userIDFromToken != userIDParam {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized access"})
			}

			c.Set("userID", userIDFromToken)
			return next(c)
		}
	}
}

func GenerateJWT(secret string, userID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("Invalid token")
}

/*
func ValidateToken(secret string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}
			claims, err := ValidateJWT(tokenParts[1], secret)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "userID", claims["sub"])
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
*/
