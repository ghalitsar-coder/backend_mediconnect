package middleware

import (
	"context"
	"net/http"
	"strings"

	"mediconnect/pkg/response"

	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Public routes should probably not use this middleware anyway
		// but we can bypass health check just in case.
		if r.URL.Path == "/api/v1/health" {
			next.ServeHTTP(w, r)
			return
		}

		// First, try reading from Authorization header (Fallback)
		var tokenStr string
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Second, try reading from HttpOnly Cookie (Primary for Web)
			cookie, err := r.Cookie("jwt_token")
			if err == nil {
				tokenStr = cookie.Value
			}
		}

		if tokenStr == "" {
			response.Error(w, http.StatusUnauthorized, "Unauthorized: missing or invalid token")
			return
		}

		// For now we'll do a simple decode just to ensure it's valid HS256/RS256 JWT
		// "supersecretkey" should ideally come from env
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("supersecretkey"), nil
		})

		if err != nil || !token.Valid {
			// If cookie is invalid, we should clear it
			http.SetCookie(w, &http.Cookie{
				Name:   "jwt_token",
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			})
			response.Error(w, http.StatusUnauthorized, "Unauthorized: invalid or expired token")
			return
		}

		// Set the user claims to the request context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
			ctx = context.WithValue(ctx, "role", claims["role"])
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			response.Error(w, http.StatusUnauthorized, "Unauthorized: invalid claims")
			return
		}
	})
}
