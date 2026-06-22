package delegator

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type DelegatorUserInfo struct {
	ID        *uuid.UUID `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
}

type DelegatorResponse struct {
	ID        *uuid.UUID         `json:"id"`
	OwnerID   uuid.UUID          `json:"owner_id"`
	Owner     *DelegatorUserInfo `json:"owner"`
	Delegate  *DelegatorUserInfo `json:"delegate"`
	StartDate time.Time          `json:"start_date"`
	EndDate   time.Time          `json:"end_date"`
	IsActive  bool               `json:"is_active"`
	CreatedAt *time.Time         `json:"created_at"`
	UpdatedAt *time.Time         `json:"updated_at"`
}
