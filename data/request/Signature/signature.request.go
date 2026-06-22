package signature

type CreateSignatureRequest struct {
	UserID   string `validate:"required,uuid" json:"userId"`
	ImageURL string `validate:"required,url" json:"imageUrl"`
}

type UpdateSignatureRequest struct {
	ImageURL string `validate:"required,url" json:"imageUrl"`
}
