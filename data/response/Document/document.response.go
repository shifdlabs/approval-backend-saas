package document

import (
	"Microservice/model"
	"time"

	uuid "github.com/satori/go.uuid"
)

type DocumentResponse struct {
	Id                  *uuid.UUID                 `json:"id"`
	Subject             string                     `json:"subject"`
	Body                string                     `json:"body"`
	Type                int                        `json:"type"`
	Step                int                        `json:"step"`
	Status              int                        `json:"status"`
	Priority            int                        `json:"priority"`
	Author              *model.User                `json:"author"`
	DocumentSequence    []model.DocumentSequence   `json:"documentSequence"`
	DocumentHistory     []model.DocumentHistory    `json:"documentHistory"`
	DocumentAttachment  []model.DocumentAttachment `json:"documentAttachment"`
	CreatedAt           time.Time                  `json:"createdAt"`
	UpdatedAt           time.Time                  `json:"updatedAt"`
	CurrentApprovalName *string                    `json:"currentApprovalName"`
	LastRejector        *RejectorResponse          `json:"lastRejector"`
}

type RejectorResponse struct {
	Name   *string `json:"name"`
	Reason *string `json:"reason"`
}
