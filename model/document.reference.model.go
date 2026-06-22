package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type DocumentReference struct {
	gorm.Model
	ID          *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	DocumentID  uuid.UUID  `gorm:"type:uuid"`
	ReferenceID uuid.UUID  `gorm:"type:uuid"`
	CreatedAt   *time.Time `gorm:"not null;default:now()"`
	UpdatedAt   *time.Time `gorm:"not null;default:now()"`
}
