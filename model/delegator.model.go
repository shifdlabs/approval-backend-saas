package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Delegator struct {
	gorm.Model
	ID         *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	OwnerID    uuid.UUID  `gorm:"type:uuid;not null"`
	DelegateID uuid.UUID  `gorm:"type:uuid;not null"`
	StartDate  time.Time  `gorm:"type:timestamp;not null"`
	EndDate    time.Time  `gorm:"type:timestamp;not null"`
	Owner      *User      `gorm:"foreignKey:OwnerID;references:ID;"`
	Delegate   *User      `gorm:"foreignKey:DelegateID;references:ID;"`
	CreatedAt  *time.Time `gorm:"not null;default:now()"`
	UpdatedAt  *time.Time `gorm:"not null;default:now()"`
}
