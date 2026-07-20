package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Position struct {
	gorm.Model
	ID             *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	OrganizationID *uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_positions_org_name"`
	Name           string     `gorm:"type:varchar(255);not null;uniqueIndex:idx_positions_org_name"`
	CreatedAt      *time.Time `gorm:"default:now()"`
	UpdatedAt      *time.Time `gorm:"not null;default:now()"`
}
