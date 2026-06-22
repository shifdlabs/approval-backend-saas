package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type CarbonCopy struct {
	gorm.Model
	ID         *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	DocumentID *uuid.UUID `gorm:"type:uuid"`
	UserID     uuid.UUID  `gorm:"type:uuid"`
	Document   *Document  `gorm:"foreignKey:DocumentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt  *time.Time `gorm:"not null;default:now()"`
}
