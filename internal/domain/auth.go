package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User represents the users table in the database
type User struct {
	ID           uuid.UUID `json:"id"`
	NIK          string    `json:"nik"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Don't expose password hash in JSON
	Phone        *string   `json:"phone,omitempty"`
	KtpURL       *string   `json:"ktp_url,omitempty"` // URL gambar KTP dari Azure Blob
	FullName     string    `json:"full_name"`
	Role         string    `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RegisterRequest represents the payload for user registration
type RegisterRequest struct {
	NIK      string  `json:"nik" validate:"required,len=16"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	Phone    *string `json:"phone,omitempty"`
	FullName string  `json:"full_name" validate:"required"`
	Role     string  `json:"role" validate:"oneof=PATIENT NAKES DINKES"` // Default: PATIENT
}

// LoginRequest represents the payload for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents the data returned after successful auth
type AuthResponse struct {
	User User `json:"user"`
	// Token is excluded here because we'll send it via HttpOnly Cookie
}

// AuthUsecase defines the interface for auth business logic
type AuthUsecase interface {
	Register(ctx context.Context, req RegisterRequest) (User, error)
	Login(ctx context.Context, req LoginRequest) (User, error) // Returns JWT token and User info
	GetUserByID(ctx context.Context, userID string) (User, error)
}

// AuthRepository defines the interface for user data operations
type AuthRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	UpdateUserKtpURL(ctx context.Context, id string, ktpURL string) error
}
