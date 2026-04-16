package usecase

import (
    "context"
    "errors"
    "fmt"

    "mediconnect/internal/domain"

    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
    authRepo  domain.AuthRepository
    jwtSecret string // tidak digunakan lagi untuk generate token, tapi mungkin untuk keperluan lain
}

func NewAuthUsecase(repo domain.AuthRepository, secret string) domain.AuthUsecase {
    return &AuthUsecase{
        authRepo:  repo,
        jwtSecret: secret,
    }
}

func (u *AuthUsecase) Register(ctx context.Context, req domain.RegisterRequest) (domain.User, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return domain.User{}, err
    }

    role := req.Role
    if role == "" {
        role = "PATIENT"
    }

    user := domain.User{
        ID:           uuid.New(),
        NIK:          req.NIK,
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
        Phone:        req.Phone,
        FullName:     req.FullName,
        Role:         role,
        IsActive:     true,
    }

    err = u.authRepo.CreateUser(ctx, &user)
    if err != nil {
        return domain.User{}, err
    }

    return user, nil
}

func (u *AuthUsecase) Login(ctx context.Context, req domain.LoginRequest) (domain.User, error) {
    user, err := u.authRepo.GetUserByEmail(ctx, req.Email)
    if err != nil {
        fmt.Println("login error:", err)
        return domain.User{}, errors.New("invalid email or password")
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
    if err != nil {
        fmt.Println("login error:", err)
        return domain.User{}, errors.New("invalid email or password")
    }

    // Tidak membuat token di sini, token dibuat di handler
    return *user, nil
}