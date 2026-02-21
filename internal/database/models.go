package database

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Uuid            string
	Email           string
	PasswordHash    string
	IsEmailVerified bool
	Status          string
}

type Roles struct {
	gorm.Model
	Uuid        string
	Name        string
	Desctiption string
}

type Permissions struct {
	gorm.Model
	Uuid        string
	Name        string
	Description string
}

type RolePermissions struct {
	RoleUuid       string
	PermissionUuid string
}

type UserRoles struct {
	UserUuid string
	RoleUuid string
}

type RefreshTokens struct {
	gorm.Model
	Uuid       string
	UserUuid   string
	TokenHash  string
	ParentUuid string
	IsRevoked  bool
	ExpiresAt  time.Time
	UserAgent  string
	IpAddress  string
}
