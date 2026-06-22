package document

type Authorize struct {
	DocumentID string `validate:"required,uuid" json:"documentId"`
	State      int    `validate:"required,oneof=1 2 3" json:"state"` // 1: approve, 2: reject, 3: cancelled
	Comment    string `json:"comment" validate:"max=1000"`
}
