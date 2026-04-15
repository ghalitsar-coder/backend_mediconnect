package handler

import (
	"mediconnect/internal/domain"
	"mediconnect/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(uc domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: uc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	token, user, err := h.authUsecase.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	maxAge := 60 * 60 * 24 // 1 hari
	secure := false        // set true saat production HTTPS
	sameSite := http.SameSiteLaxMode

	// 1. HttpOnly JWT — tidak terbaca JS/browser, aman dari XSS
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	})

	// 2. auth_session flag — non-HttpOnly, dibaca oleh proxy.ts di Edge
	//    Tidak mengandung data sensitif, hanya penanda "sudah login"
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_session",
		Value:    "1",
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: false,
		Secure:   secure,
		SameSite: sameSite,
	})

	// 3. user_role — non-HttpOnly, dibaca oleh proxy.ts untuk routing per-role
	//    Tidak sensitif karena role bersifat publik di konteks routing
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

// Logout — clear semua auth cookies
func (h *AuthHandler) Logout(c *gin.Context) {
	expired := time.Unix(0, 0)

	cookiesToClear := []string{"token", "auth_session", "user_role"}
	for _, name := range cookiesToClear {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			Expires:  expired,
			MaxAge:   -1,
			HttpOnly: name == "token", // token tetap HttpOnly
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
	}

	response.Success(c, http.StatusOK, "Logged out successfully", nil)
}
