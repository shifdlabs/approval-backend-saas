package user

type UpdateEmailRequest struct {
	NewEmail string `validate:"required,email" json:"newEmail"`
}
