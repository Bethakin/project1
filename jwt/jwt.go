package utils

import (
	"fmt"
	"net/http"
	"strings"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type contextKey string

const userIDKey contextKey = "userID"

func GetUserIDFromContext(ctx context.Context) (int, bool) {
	val := ctx.Value(userIDKey)
	userID, ok := val.(int)
	return userID, ok
}

/*
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
}*/

// AuthMiddleware parses the JWT and injects userID into context
func AuthMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization format")
			}

			claims, err := ValidateJWT(tokenString, secret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token: "+err.Error())
			}

			sub, ok := claims["sub"]
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token missing 'sub' claim")
			}

			var userID int
			switch v := sub.(type) {
			case float64:
				userID = int(v)
			case int:
				userID = v
			default:
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid 'sub' type in token")
			}

			// Inject userID into context
			ctx := context.WithValue(c.Request().Context(), userIDKey, userID)
			c.SetRequest(c.Request().WithContext(ctx))


			return next(c)
		}
	}
}

func GenerateJWT(secret string, userID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID, // should be an integer
		"exp": time.Now().Add(24 * time.Hour).Unix(),
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
