package user

type DeleteMultipleUserRequest struct {
	IDs []string `validate:"required,min=1,max=100,dive,required" json:"ids"`
}
