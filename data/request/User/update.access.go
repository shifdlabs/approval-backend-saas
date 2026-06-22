package user

type UpdateAccessRequest struct {
	ID     string `validate:"required,uuid" json:"id"`
	Access bool   `json:"access"`
}
