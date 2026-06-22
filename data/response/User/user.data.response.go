package response

import (
	"Microservice/model"
	"time"

	uuid "github.com/satori/go.uuid"
)

type UserResponse struct {
	ID        *uuid.UUID      `json:"id"`
	Email     string          `json:"email"`
	Role      int             `json:"role"`
	FirstName string          `json:"firstName"`
	LastName  string          `json:"lastName"`
	Position  *model.Position `json:"position"`
	Access    bool            `json:"access"`
	Phone     string          `json:"phone"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}
