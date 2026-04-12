package handler

import (
	"mediconnect/internal/domain"
	"mediconnect/pkg/response"
	"net/http"

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
