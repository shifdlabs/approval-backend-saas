package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Recipient struct {
	gorm.Model
	DocumentID uuid.UUID  `gorm:"type:uuid"`
	UserID     uuid.UUID  `gorm:"type:uuid"`
	CreatedAt  *time.Time `gorm:"not null;default:now()"`
	UpdatedAt  *time.Time `gorm:"not null;default:now()"`
	Document   *Document  `gorm:"foreignKey:DocumentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
