package user

type UpdateBiodataRequest struct {
	PositionID string `json:"position" validate:"omitempty,uuid"`
	FirstName  string `validate:"required,min=1,max=200" json:"firstName"`
	LastName   string `validate:"required,min=1,max=200" json:"lastName"`
	Phone      string `validate:"required,min=5,max=20" json:"phone"`
}
