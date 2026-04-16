package handler

import (
    "mediconnect/internal/domain"
    "mediconnect/pkg/jwt"
    "mediconnect/pkg/response"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authUsecase domain.AuthUsecase
    jwtManager  *jwt.JWTManager
}

func NewAuthHandler(uc domain.AuthUsecase, jwtManager *jwt.JWTManager) *AuthHandler {
    return &AuthHandler{
        authUsecase: uc,
        jwtManager:  jwtManager,
    }
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req domain.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request format")
        return
    }

    user, err := h.authUsecase.Login(c.Request.Context(), req)
    if err != nil {
        response.Error(c, http.StatusUnauthorized, "Invalid credentials")
        return
    }

    // Generate token dengan JWTManager
    token, err := h.jwtManager.GenerateAccessToken(user.ID.String(), user.Email, user.Role)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to generate token")
        return
    }

    maxAge := 60 * 60 * 24 // 1 hari
    secure := false        // true saat production HTTPS
    sameSite := http.SameSiteLaxMode

    // 1. HttpOnly JWT
    http.SetCookie(c.Writer, &http.Cookie{
        Name:     "token",
        Value:    token,
        Path:     "/",
        MaxAge:   maxAge,
        HttpOnly: true,
        Secure:   secure,
        SameSite: sameSite,
    })

    // 2. auth_session flag
    http.SetCookie(c.Writer, &http.Cookie{
        Name:     "auth_session",
        Value:    "1",
        Path:     "/",
        MaxAge:   maxAge,
        HttpOnly: false,
        Secure:   secure,
        SameSite: sameSite,
    })

    // 3. user_role
    http.SetCookie(c.Writer, &http.Cookie{
        Name:     "user_role",
        Value:    user.Role,
        Path:     "/",
        MaxAge:   maxAge,
        HttpOnly: false,
        Secure:   secure,
        SameSite: sameSite,
    })

    response.Success(c, http.StatusOK, "Login successful", map[string]interface{}{
        "token": token,
        "user":  user,
    })
}

// Register dan Logout tetap sama seperti kode sebelumnya
func (h *AuthHandler) Register(c *gin.Context) {
    var req domain.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request format")
        return
    }

    _, err := h.authUsecase.Register(c.Request.Context(), req)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to register user")
        return
    }

    response.Success(c, http.StatusCreated, "User registered successfully", nil)
}

func (h *AuthHandler) Logout(c *gin.Context) {
    expired := time.Unix(0, 0)
    cookiesToClear := []struct {
        name     string
        httpOnly bool
    }{
        {"token", true},
        {"auth_session", false},
        {"user_role", false},
    }

    for _, cookie := range cookiesToClear {
        http.SetCookie(c.Writer, &http.Cookie{
            Name:     cookie.name,
            Value:    "",
            Path:     "/",
            Expires:  expired,
            MaxAge:   -1,
            HttpOnly: cookie.httpOnly,
            Secure:   false,
            SameSite: http.SameSiteLaxMode,
        })
    }

    response.Success(c, http.StatusOK, "Logged out successfully", nil)
}