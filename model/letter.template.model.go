package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type LetterTemplate struct {
	gorm.Model
	ID             uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	OrganizationID *uuid.UUID `gorm:"type:uuid;not null"`
	Name           string     `gorm:"type:varchar(255);not null"`
	Description    string     `gorm:"type:varchar(255)"`
	Body           string     `gorm:"type:text"`
	CreatedAt      *time.Time `gorm:"default:now()"`
	UpdatedAt      *time.Time `gorm:"not null;default:now()"`
}
