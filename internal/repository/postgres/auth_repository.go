package postgres

import (
	"context"
	"mediconnect/internal/domain"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) domain.AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Table("users").Create(user).Error
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Table("users").Where("email = ? AND is_active = ?", email, true).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Table("users").Where("id = ? AND is_active = ?", id, true).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) UpdateUserKtpURL(ctx context.Context, id string, ktpURL string) error {
	return r.db.WithContext(ctx).Table("users").Where("id = ?", id).Update("ktp_url", ktpURL).Error
}
