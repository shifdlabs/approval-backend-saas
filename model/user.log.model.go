package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UserLog struct {
	gorm.Model
	ID             *uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	OrganizationID *uuid.UUID     `gorm:"type:uuid;not null"`
	UserID         uuid.UUID      `gorm:"type:uuid"`
	Action         string         `gorm:"type:varchar"`
	Module         string         `gorm:"type:varchar"`
	Log            datatypes.JSON `gorm:"type:jsonb"`
	LogDate        *time.Time     `gorm:"not null;default:now()"`
}
