package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type AppSettings struct {
	gorm.Model
	OrganizationID *uuid.UUID `gorm:"type:uuid;not null"`
	Key            string     `gorm:"type:varchar"`
	Value          string     `gorm:"type:varchar"`
	CreatedAt      *time.Time `gorm:"not null;default:now()"`
	UpdatedAt      *time.Time `gorm:"not null;default:now()"`
}
