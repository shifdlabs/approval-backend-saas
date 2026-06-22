package attachment

type AttachmentRequest struct {
	Id           string `validate:"required" json:"id"`
	FileName     string `validate:"required" json:"fileName"`
	OriginalName string `validate:"required" json:"originalName"`
}
