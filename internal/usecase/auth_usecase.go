package usecase

import (
	"context"
	"errors"
	"time"

	"mediconnect/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	authRepo  domain.AuthRepository
	jwtSecret string
}

func NewAuthUsecase(repo domain.AuthRepository, secret string) domain.AuthUsecase {
	return &AuthUsecase{
		authRepo:  repo,
		jwtSecret: secret,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, req domain.RegisterRequest) (domain.User, error) {
	// 1. Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}

	role := req.Role
	if role == "" {
		role = "PATIENT"
	}

	// 2. Create the user model
	user := domain.User{
		NIK:          req.NIK,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Phone:        req.Phone,
		FullName:     req.FullName,
		Role:         role,
		IsActive:     true,
	}

	// 3. Save to database
	err = u.authRepo.CreateUser(ctx, &user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (u *AuthUsecase) Login(ctx context.Context, req domain.LoginRequest) (string, domain.User, error) {
	// 1. Find user by email
	user, err := u.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", domain.User{}, errors.New("invalid email or password")
	}

	// 2. Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", domain.User{}, errors.New("invalid email or password")
	}

	// 3. Generate JWT
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", domain.User{}, err
	}

	return tokenString, *user, nil
}
