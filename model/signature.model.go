package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Signature struct {
	gorm.Model
	ID        *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID    *uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	ImageURL  string     `gorm:"type:text;not null"`
	CreatedAt *time.Time `gorm:"default:now()"`
	UpdatedAt *time.Time `gorm:"not null;default:now()"`
	User      *User      `gorm:"foreignKey:UserID;references:ID"`
}
