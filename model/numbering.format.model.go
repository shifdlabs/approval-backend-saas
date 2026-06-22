package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type NumberingFormat struct {
	gorm.Model
	ID               *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name             string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	GroupID          uuid.UUID      `gorm:"type:uuid"`
	Format           string         `gorm:"type:varchar(255);not null"`
	Separator        string         `gorm:"type:varchar(3);not null;default:'/'"`
	IncrementByGroup *bool          `gorm:"type:bool;not null;"`
	CreatedAt        *time.Time     `gorm:"default:now()"`
	UpdatedAt        *time.Time     `gorm:"not null;default:now()"`
	Group            NumberingGroup `gorm:"foreignKey:GroupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
