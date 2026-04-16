package middleware

import (
	"log"
	"mediconnect/pkg/jwt"
	"mediconnect/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth(jwtManager *jwt.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("JWTAuth middleware triggered for:", c.Request.URL.Path)
		if c.Request.URL.Path == "/api/v1/health" {
			c.Next()
			return
		}

		var tokenStr string

		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			cookie, err := c.Cookie("token")
			if err == nil {
				tokenStr = cookie
			} else {
				cookie, err = c.Cookie("jwt_token")
				if err == nil {
					tokenStr = cookie
				}
			}
		}

		if tokenStr == "" {
			response.Error(c, http.StatusUnauthorized, "Unauthorized: missing token")
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateToken(tokenStr)
		if err != nil {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.SetCookie("jwt_token", "", -1, "/", "", false, true)
			response.Error(c, http.StatusUnauthorized, "Unauthorized: "+err.Error())
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
