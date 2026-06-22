package documenthistory

import (
	document "Microservice/data/response/Document"

	uuid "github.com/satori/go.uuid"
)

type DocumentHistoryResponse struct {
	Id          *uuid.UUID                `json:"id"`
	Step        int                       `json:"step"`
	DocumentID  *uuid.UUID                `json:"documentId"`
	UserID      *uuid.UUID                `json:"userId"`
	Description string                    `json:"description"`
	IsApproved  bool                      `json:"is_approved"`
	Document    document.DocumentResponse `json:"document"`
}
