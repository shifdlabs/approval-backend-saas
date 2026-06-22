package document

import (
	"Microservice/model"
	"time"

	uuid "github.com/satori/go.uuid"
)

type EditDocumentResponse struct {
	Id                    *uuid.UUID            `json:"id"`
	PublicationNumberType int                   `json:"publicationNumberType"` // 1: Auto-Generated, 2: Booking Number, 3: Custom, 4: N/A (No Number)
	PublicationValue      *string               `json:"publicationValue"`      // it could be Booked Number, Format ID or Custom Number
	Subject               string                `json:"subject"`
	Body                  string                `json:"body"`
	Type                  int                   `json:"type"`
	Step                  int                   `json:"step"`
	Status                int                   `json:"status"`
	Priority              int                   `json:"priority"`
	Author                model.User            `json:"author"`
	DocumentAttachment    *[]DocumentAttachment `json:"documentAttachment"`
	DocumentReferences    *[]DocumentReference  `json:"documentReferences"`
	ExternalRecipient     *string               `json:"externalRecipient"`
	InternalRecipients    *[]string             `json:"internalRecipients"`
	CarbonCopy            *[]string             `json:"carbonCopy"`
	Approvers             *[]string             `json:"approvers"`
	CreatedAt             time.Time             `json:"createdAt"`
	UpdatedAt             time.Time             `json:"updatedAt"`
}
