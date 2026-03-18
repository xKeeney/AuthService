package auth

import (
	"errors"
)

var (
	ErrUserExist           = errors.New("user already exist")
	ErrUserNotFound        = errors.New("user not found")
	ErrWrongPassword       = errors.New("wrong password")
	ErrAccessTokenExp      = errors.New("access token expired")
	ErrRefreshTokenExp     = errors.New("refresh token expired")
	ErrInvalidAccessToken  = errors.New("invalid access token")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrParseJWT            = errors.New("cannot parse claims")
	ErrLoadPrivateKey      = errors.New("cannot load private key")
	ErrLoadPublicKey       = errors.New("cannot load public key")
)
