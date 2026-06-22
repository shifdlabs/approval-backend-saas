package bookmark

type BookmarkRequest struct {
	UserID     string `validate:"required,uuid" json:"userId"`     // ID pengguna yang login
	DocumentID string `validate:"required,uuid" json:"documentId"` // ID dokumen yang di-bookmark
}
