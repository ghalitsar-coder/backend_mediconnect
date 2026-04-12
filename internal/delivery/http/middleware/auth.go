package middleware

import (
  "net/http"
  "strings"
  "mediconnect/pkg/response"
  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v5"
)

func JWTAuth(c *gin.Context) {
  if c.Request.URL.Path == "/api/v1/health" {
      c.Next()
      return
  }

  var tokenStr string
  authHeader := c.GetHeader("Authorization")
  if strings.HasPrefix(authHeader, "Bearer ") {
      tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
  } else {
      cookie, err := c.Cookie("jwt_token")
      if err == nil {
          tokenStr = cookie
      }
  }

  if tokenStr == "" {
      response.Error(c, http.StatusUnauthorized, "Unauthorized: missing token")
      c.Abort()
      return
  }

  token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
      return []byte("supersecretkey"), nil
  })

  if err != nil || !token.Valid {
      c.SetCookie("jwt_token", "", -1, "/", "", false, true)
      response.Error(c, http.StatusUnauthorized, "Unauthorized: invalid/expired token")
      c.Abort()
      return
  }

  if claims, ok := token.Claims.(jwt.MapClaims); ok {
      c.Set("user_id", claims["user_id"])
      c.Set("role", claims["role"])
      c.Next()
  } else {
      response.Error(c, http.StatusUnauthorized, "Unauthorized: invalid claims")
      c.Abort()
      return
  }
}