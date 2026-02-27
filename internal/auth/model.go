package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTCustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}
