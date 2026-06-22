package documentSequence

import (
	uuid "github.com/satori/go.uuid"
)

type DocumentAttachmentResponse struct {
	Id           *uuid.UUID `json:"id"`
	DocumentID   *uuid.UUID `json:"documentId"`
	OriginalName string     `json:"originalName"`
	FileName     string     `json:"fileName"`
	Path         string     `json:"path"`
	Size         string     `json:"size"`
	Type         string     `json:"type"`
}
