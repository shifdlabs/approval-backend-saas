package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Bookmark struct {
	BookmarkID uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null"`
	DocumentID uuid.UUID  `gorm:"type:uuid;not null"`
	CreatedAt  *time.Time `gorm:"not null;default:now()"`
	Document   *Document  `gorm:"foreignKey:DocumentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
