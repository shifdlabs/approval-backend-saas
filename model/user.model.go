package model

import (
	"html"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID             *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	OrganizationID *uuid.UUID `gorm:"type:uuid;not null"`
	PositionID     *uuid.UUID `gorm:"type:uuid"`
	EmployeeID string     `gorm:"type:varchar(100)"`
	Email      string     `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password   string     `gorm:"size:255;not null" json:"-"`
	Role          int        `gorm:"type:integer"` // 99: admin, 1: user
	FirstName     string     `gorm:"type:varchar(100)"`
	LastName      string     `gorm:"type:varchar(100)"`
	Access        bool       `gorm:"type:boolean"`
	IsLocked      bool       `gorm:"type:boolean;default:false" json:"-"`
	LockTimestamp *time.Time `gorm:"type:timestamp" json:"-"`
	Phone         string     `gorm:"type:varchar(20)"`
	Position      *Position  `gorm:"foreignKey:PositionID;constraint:OnDelete:SET NULL;"`
	CreatedAt  *time.Time `gorm:"not null;default:now()"`
	UpdatedAt  *time.Time `gorm:"not null;default:now()"`
	DeletedAt  *time.Time
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	u.Email = html.EscapeString(strings.TrimSpace(u.Email))

	return nil
}

func HashPasswordString(password string) (*string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashed := string(hashedPassword)

	return &hashed, nil
}
