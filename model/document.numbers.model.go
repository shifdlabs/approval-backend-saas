package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type DocumentNumbers struct {
	gorm.Model
	ID                *uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Value             string           `gorm:"type:varchar(255);not null"`
	DocumentID        *uuid.UUID       `gorm:"type:uuid"`
	NumberingFormatID *uuid.UUID       `gorm:"type:uuid"`
	UserId            *string          `gorm:"type:varchar(255)"`
	State             int              `gorm:"type:int"` // 1: Booked, 2: Saved, 0: Cancelled
	CreatedAt         *time.Time       `gorm:"default:now()"`
	UpdatedAt         *time.Time       `gorm:"not null;default:now()"`
	NumberingFormat   *NumberingFormat `gorm:"foreignKey:NumberingFormatID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Document          *Document        `gorm:"foreignKey:DocumentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
