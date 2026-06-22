package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type FailedLoginAttempt struct {
	gorm.Model
	ID          *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID      *uuid.UUID `gorm:"type:uuid;not null"`
	User        *User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	AttemptedAt *time.Time `gorm:"not null;default:now()"`
	CreatedAt   *time.Time `gorm:"not null;default:now()"`
	UpdatedAt   *time.Time `gorm:"not null;default:now()"`
	DeletedAt   *time.Time
}
