package user

type CreateUserRequest struct {
	PositionID string `json:"positionID" validate:"omitempty,uuid"`
	EmployeeID string `json:"employeeID"`
	Email      string `validate:"required,email" json:"email"`
	Password   string `validate:"required,min=8,max=200" json:"password"`
	Role       int    `validate:"required,oneof=1 99" json:"role"`
	FirstName  string `validate:"required,min=1,max=200" json:"firstName"`
	LastName   string `validate:"required,min=1,max=200" json:"lastName"`
	Access     bool   `json:"access"`
	Phone      string `validate:"required,min=5,max=20" json:"phone"`
}
