package database

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	UUID string `gorm:"type:uuid;primaryKey"`

	Email           string `gorm:"not null;uniqueIndex"`
	PasswordHash    string `gorm:"not null"`
	IsEmailVerified bool   `gorm:"not null;default:false"`

	// active | blocked | deleted
	Status string `gorm:"type:varchar(16);not null;index"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Roles []Role `gorm:"many2many:user_roles"`
}

type Role struct {
	UUID string `gorm:"type:uuid;primaryKey"`

	Name        string `gorm:"type:varchar(64);not null;uniqueIndex"`
	Description string

	Permissions []Permission `gorm:"many2many:role_permissions"`
}

type Permission struct {
	UUID string `gorm:"type:uuid;primaryKey"`

	Name        string `gorm:"type:varchar(128);not null;uniqueIndex"`
	Description string
}

type UserRole struct {
	UserUUID string `gorm:"type:uuid;primaryKey"`
	RoleUUID string `gorm:"type:uuid;primaryKey"`
}

type RolePermission struct {
	RoleUUID       string `gorm:"type:uuid;primaryKey"`
	PermissionUUID string `gorm:"type:uuid;primaryKey"`
}

type RefreshToken struct {
	UUID string `gorm:"type:uuid;primaryKey"`

	UserUUID string `gorm:"type:uuid;not null;index"`
	User     User   `gorm:"foreignKey:UserUUID;references:UUID"`

	TokenHash string `gorm:"not null;uniqueIndex"`

	ParentUUID *string `gorm:"type:uuid;index"`

	IsRevoked bool      `gorm:"not null;default:false;index"`
	ExpiresAt time.Time `gorm:"not null;index"`

	CreatedAt time.Time
}

type AuditLog struct {
	UUID string `gorm:"type:uuid;primaryKey"`

	UserUUID *string `gorm:"type:uuid;index"`

	Action string `gorm:"type:varchar(32);not null;index"`

	IPAddress string `gorm:"type:inet"`
	UserAgent string
	CreatedAt time.Time
}
