package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type DocumentAttachment struct {
	gorm.Model
	ID           *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	DocumentID   uuid.UUID  `gorm:"type:uuid"`
	OriginalName string     `gorm:"type:varchar"`
	FileName     string     `gorm:"type:varchar"`
	Path         string     `gorm:"type:varchar"`
	Size         string     `gorm:"type:varchar"`
	Type         string     `gorm:"type:varchar"`
	Document     *Document  `gorm:"foreignKey:DocumentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt    *time.Time `gorm:"not null;default:now()"`
	UpdatedAt    *time.Time `gorm:"not null;default:now()"`
}
