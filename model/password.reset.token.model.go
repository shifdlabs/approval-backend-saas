package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type PasswordResetToken struct {
	gorm.Model
	ID        *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID    *uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenHash string     `gorm:"type:varchar(64);not null;uniqueIndex"`
	ExpiresAt *time.Time `gorm:"not null"`
	UsedAt    *time.Time
	CreatedAt *time.Time `gorm:"not null;default:now()"`
	UpdatedAt *time.Time `gorm:"not null;default:now()"`
	DeletedAt *time.Time
	User      *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}
