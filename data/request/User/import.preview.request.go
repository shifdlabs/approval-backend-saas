package user

type PreviewImportRequest struct {
	ColumnMapping map[string]string `json:"columnMapping"`
}

type ImportedUserData struct {
	EmployeeID string `json:"employeeID"`
	Email      string `json:"email" validate:"required,email"`
	FirstName  string `json:"firstName" validate:"required,min=1,max=200"`
	LastName   string `json:"lastName" validate:"required,min=1,max=200"`
	Role       int    `json:"role" validate:"omitempty,oneof=1 99"`
	Phone      string `json:"phone" validate:"omitempty,min=5,max=20"`
	PositionID string `json:"positionID" validate:"omitempty,uuid"`
	Password   string `json:"password"`
}
