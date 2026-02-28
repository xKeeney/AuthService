package auth

import (
	"auth_service/internal/database"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func InitAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) CreateUser(uuid, email, passwordHash, status string) error {
	user := database.User{
		UUID:         uuid,
		Email:        email,
		PasswordHash: passwordHash,
		Status:       status,
	}
	if err := r.db.Create(&user).Error; err != nil {
		return fmt.Errorf("create_user(%s) error: %v", email, err)
	}
	return nil
}

func (r *authRepository) SelectUserByEmail(email string) (*database.User, error) {
	var user database.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("select_user_by_email(%s) error: %v", email, err)
		}
	}
	return &user, nil
}

func (r *authRepository) CreateRefreshToken(uuid, userUUID, tokenHash string, parentUUID *string, expiresAt time.Time) error {
	refresh := database.RefreshToken{
		UUID:       uuid,
		UserUUID:   userUUID,
		TokenHash:  tokenHash,
		ParentUUID: parentUUID,
		ExpiresAt:  expiresAt,
	}
	if err := r.db.Create(&refresh).Error; err != nil {
		return fmt.Errorf("create_refresh_token(userUUID='%s') error: %v", userUUID, err)
	}
	return nil
}
