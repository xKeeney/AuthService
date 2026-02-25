package auth

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	authRepo *authRepository
}

func InitAuthService(authRepo *authRepository) *authService {
	return &authService{
		authRepo: authRepo,
	}
}

func (s *authService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash_password error: %v", err)
	}
	return string(bytes), nil
}

func (s *authService) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *authService) RegisterUser(email, password string) (bool, error) {
	// Select user by email
	user, err := s.authRepo.SelectUserByEmail(email)
	if err != nil {
		return false, fmt.Errorf("register_user(%s) error: %v", email, err)
	}

	// Check if user exist
	if user != nil {
		return false, nil
	}

	// User data
	uuidStr := uuid.NewString()
	passHash, err := s.hashPassword(password)
	if err != nil {
		return false, fmt.Errorf("register_user(%s) error: %v", email, err)
	}
	status := "active"

	// Create user
	if err := s.authRepo.CreateUser(uuidStr, email, passHash, status); err != nil {
		return false, fmt.Errorf("register_user(%s) error: %v", email, err)
	}

	return true, nil
}

func (s *authService) CreateUser(uuid, email, passwordHash, status string) error {
	return s.authRepo.CreateUser(uuid, email, passwordHash, status)
}
