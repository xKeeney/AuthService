package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

/* ADDITIONAL LOGIC */

// Create hash for password
func (s *authService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash_password error: %v", err)
	}
	return string(bytes), nil
}

// Validate password by hash
func (s *authService) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Read RSA private key from file
func (s *authService) loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("invalid private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// Read RSA public key from file
func (s *authService) loadPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("invalid public key")
	}

	return x509.ParsePKCS1PublicKey(block.Bytes)
}

// Generate JWT (access_token)
func (s *authService) generateJWT(privateKey *rsa.PrivateKey, userUuid string, ttl time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userUuid,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "AuthService",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(privateKey)
}

// Validate JWT (access token)
func (s *authService) validateJWT(tokenString string, publicKey *rsa.PublicKey) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// Get expired JWT claims
func (s *authService) parseExpiredJWT(tokenString string, publicKey *rsa.PublicKey) (*jwt.RegisteredClaims, error) {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	token, err := parser.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func (s *authService) generateRefreshToken(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	result.Grow(length)

	for _, v := range b {
		result.WriteByte(charset[v%byte(len(charset))])
	}

	return result.String(), nil
}

func (s *authService) hashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *authService) checkRefreshTokenHash(token, hash string) bool {
	tmpHash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(tmpHash[:])
	if tokenHash == hash {
		return true
	} else {
		return false
	}
}

/* MAIN LOGIC */

// Register
func (s *authService) RegisterUser(email, password string) error {
	// Select user by email
	user, err := s.authRepo.SelectUserByEmail(email)
	if err != nil {
		return fmt.Errorf("register_user(%s) error: %v", email, err)
	}

	// Check if user exist
	if user != nil {
		return ErrUserExist
	}

	// User data
	uuidStr := uuid.NewString()
	passHash, err := s.hashPassword(password)
	if err != nil {
		return fmt.Errorf("register_user(%s) error: %v", email, err)
	}
	status := "active"

	// Create user
	if err := s.authRepo.CreateUser(uuidStr, email, passHash, status); err != nil {
		return fmt.Errorf("register_user(%s) error: %v", email, err)
	}

	return nil
}

func (s *authService) LoginUser(email, password string) (string, string, error) {
	// Select user by email
	user, err := s.authRepo.SelectUserByEmail(email)
	if err != nil {
		return "", "", fmt.Errorf("login_user(%s) error: %v", email, err)
	}

	// Check is user exist
	if user == nil {
		return "", "", ErrUserNotFound
	}

	// Check is password correct
	if ok := s.checkPasswordHash(password, user.PasswordHash); !ok {
		return "", "", ErrWrongPassword
	}

	// load jwt private key
	privKey, err := s.loadPrivateKey("private.pem")
	if err != nil {
		return "", "", fmt.Errorf("login_user(%s) error: %v", email, err)
	}

	// create access_token
	JWTTtl := 15 * time.Minute
	accessToken, err := s.generateJWT(privKey, user.UUID, JWTTtl)
	if err != nil {
		return "", "", fmt.Errorf("login_user(%s) error: %v", email, err)
	}

	// generate refresh token
	refreshToken, err := s.generateRefreshToken(256)
	if err != nil {
		return "", "", fmt.Errorf("login_user(%s) error: %v", email, err)
	}

	// hash refresh token
	refreshTokenHash := s.hashRefreshToken(refreshToken)

	// add refresh token to db
	refreshTokenUUID := uuid.NewString()
	refreshTokenTTL := 360 * time.Hour
	expireTime := time.Now().Add(refreshTokenTTL)
	if err := s.authRepo.CreateRefreshToken(refreshTokenUUID, user.UUID, refreshTokenHash, nil, expireTime); err != nil {
		return "", "", fmt.Errorf("login_user(%s) error: %v", email, err)
	}

	return accessToken, refreshToken, nil
}
