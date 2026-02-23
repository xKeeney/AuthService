package auth

import (
	"auth_service/internal/database"

	"gorm.io/gorm"
)

type authRepo interface {
	CreateUser(uuid, email, passwordHash, status string) error
}

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
	return r.db.Create(&user).Error
}
