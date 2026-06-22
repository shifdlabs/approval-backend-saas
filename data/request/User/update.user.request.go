package user

type UpdateUserRequest struct {
	ID         string `validate:"required,uuid" json:"id"`
	PositionID string `json:"position" validate:"omitempty,uuid"`
	EmployeeID string `json:"employeeID"`
	Email      string `validate:"required,email" json:"email"`
	Role       int    `validate:"required,oneof=1 99" json:"role"`
	FirstName  string `validate:"required,min=1,max=200" json:"firstName"`
	LastName   string `validate:"required,min=1,max=200" json:"lastName"`
	Access     bool   `json:"access"`
	Phone      string `validate:"required,min=5,max=20" json:"phone"`
}

type UpdateUserTypeRequest struct {
	ID   string `validate:"required,uuid" json:"id"`
	Type int    `validate:"required" json:"type"`
}
