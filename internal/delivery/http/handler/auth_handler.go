package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"mediconnect/internal/domain"
	"mediconnect/pkg/response"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(uc domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: uc}
}

// Register godoc
//
//	@Summary      Register a new user
//	@Tags         auth
//	@Accept       json
//	@Produce      json
//	@Param        request  body  domain.RegisterRequest  true  "Register Payload"
//	@Success      201  {object}  response.APIResponse
//	@Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.authUsecase.Register(r.Context(), req)
	if err != nil {
		// Example: duplicate entity, etc.
		response.Error(w, http.StatusInternalServerError, "Failed to register user: "+err.Error())
		return
	}

	response.Success(w, http.StatusCreated, "User registered successfully", user)
}

// Login godoc
//
//	@Summary      Login and receive HttpOnly JWT
//	@Tags         auth
//	@Accept       json
//	@Produce      json
//	@Param        request  body  domain.LoginRequest  true  "Login Payload"
//	@Success      200  {object}  response.APIResponse
//	@Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	token, user, err := h.authUsecase.Login(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Set JWT as HttpOnly Cookie

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // Set to false for LOCAL testing (HTTP). Change to true in Prod.
		SameSite: http.SameSiteLaxMode,
	})

	response.Success(w, http.StatusOK, "Login successful", domain.AuthResponse{User: user})
}

// Logout clears the HttpOnly JWT Cookie
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Erase the cookie by setting it to expire in the past
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	response.Success(w, http.StatusOK, "Logout successful", nil)
}
