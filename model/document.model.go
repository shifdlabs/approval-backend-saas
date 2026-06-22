package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	ID                      uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	OrganizationID          *uuid.UUID `gorm:"type:uuid;not null"`
	AuthorID                *uuid.UUID `gorm:"type:uuid"`
	PublicationNumberType   int        `gorm:"type:int;not null"` // 1: Auto-Generatedm 2: Booking Number, 3: Custom, 4: N/A (No Number)
	CustomPublicationNumber *string    `gorm:"type:varchar"`
	Type                    int        `gorm:"type:int"`
	Priority                int        `gorm:"type:int"` // 1: High, 2: Medium, 3: Low
	Subject                 string     `gorm:"type:varchar"`
	Body                    string     `gorm:"type:varchar"`
	ExternalRecipient       string     `gorm:"type:varchar"`
	Step                    int        `gorm:"type:int"`
	LetterHead              bool       `gorm:"type:bool;not null"`
	TemplateID              *uuid.UUID `gorm:"type:uuid"`
	Status                  int        `gorm:"type:int;not null"` // 0: draft, 1: Inprogress, 2: finished, 3: Cancelled, 99: rejected
	Author                  *User      `gorm:"foreignKey:AuthorID;constraint:OnDelete:SET NULL;"`
	DocumentSequence        []DocumentSequence
	DocumentHistory         []DocumentHistory
	DocumentAttachment      []DocumentAttachment
	Recipients              []Recipient
	DueDate                 *time.Time `gorm:"type:timestamp"`
	CreatedAt               *time.Time `gorm:"default:now()"`
	UpdatedAt               *time.Time `gorm:"not null;default:now()"`
}
