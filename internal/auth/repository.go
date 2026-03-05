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

func (r *authRepository) DeleteRefreshTokenByTokenHash(tokenHash string) error {
	if err := r.db.Where("token_hash = ?", tokenHash).Delete(&database.RefreshToken{}).Error; err != nil {
		return fmt.Errorf("delete_refresh_token_by_token_hash(%s) error: %v", tokenHash, err)
	}
	return nil
}

func (r *authRepository) SelectActiveRefreshTokensByUserUUID(userUUID string) ([]database.RefreshToken, error) {
	var refreshTokens []database.RefreshToken
	if err := r.db.Where(&database.RefreshToken{IsRevoked: false, ExpiresAt: time.Now(), UserUUID: userUUID}).Find(&refreshTokens).Error; err != nil {
		return nil, fmt.Errorf("select_active_refresh_tokens_by_user_uuid(%s) error: %v", userUUID, err)
	}
	return refreshTokens, nil
}

func (r *authRepository) SelectRefreshTokenByTokenHash(tokenHash string) (*database.RefreshToken, error) {
	var refreshToken database.RefreshToken
	if err := r.db.Where(&database.RefreshToken{TokenHash: tokenHash}).Find(&refreshToken).Error; err != nil {
		return nil, fmt.Errorf("select_refresh_tokens_by_token_hash error: %v", err)
	}
	return &refreshToken, nil
}

func (r *authRepository) UpdateRefreshToken(newRefreshToken database.RefreshToken) error {
	if err := r.db.Save(&newRefreshToken).Error; err != nil {
		return fmt.Errorf("update_refresh_token(tokenUUID=%s) error: %v", newRefreshToken.UUID, err)
	}
	return nil
}
