package user

type UpdatePasswordRequest struct {
	ID              string `validate:"required,uuid" json:"id"`
	CurrentPassword string `validate:"omitempty,min=8,max=200" json:"currentPassword"`
	NewPassword     string `validate:"required,min=8,max=200" json:"newPassword"`
}
