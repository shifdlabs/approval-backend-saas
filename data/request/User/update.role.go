package user

type UpdateRoleRequest struct {
	ID   string `validate:"required,uuid" json:"id"`
	Role int    `validate:"required,oneof=1 99" json:"role"`
}
