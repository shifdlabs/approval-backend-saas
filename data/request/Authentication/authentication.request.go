package authentication

import (
	"Microservice/model"
	"time"

	uuid "github.com/satori/go.uuid"
)

type VerifyForgetPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPassword struct {
	Email       string `json:"email" validate:"required,email"`
	NewPassword string `json:"newPassword" validate:"required,min=8,max=200"`
}

type UserResponse struct {
	ID        *uuid.UUID `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Email     string     `json:"email,omitempty"`
	Role      string     `json:"role,omitempty"`
	Photo     string     `json:"photo,omitempty"`
	Provider  string     `json:"provider"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func FilterUserRecord(user *model.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: *user.CreatedAt,
		UpdatedAt: *user.UpdatedAt,
	}
}
