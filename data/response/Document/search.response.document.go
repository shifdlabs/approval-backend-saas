package document

import uuid "github.com/satori/go.uuid"

type SearchDocumentResponse struct {
	ID             uuid.UUID `json:"id"`
	Subject        string    `json:"subject"`
	DocumentNumber *string   `json:"document_number"` // null jika type 4
}
