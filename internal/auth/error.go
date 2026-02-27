package auth

import (
	"errors"
)

var (
	ErrUserExist       = errors.New("user already exist")
	ErrUserNotFound    = errors.New("user not found")
	ErrWrongPassword   = errors.New("wrong password")
	ErrAccessTokenExp  = errors.New("access token expired")
	ErrRefreshTokenExp = errors.New("refresh token expired")
	ErrInvalidToken    = errors.New("invalid token")
	ErrParseJWT        = errors.New("cannot parse claims")
)
