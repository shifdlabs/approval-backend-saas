package document

import (
	"Microservice/model"
	"time"

	uuid "github.com/satori/go.uuid"
)

type DocumentDetailResponse struct {
	Id                 *uuid.UUID                 `json:"id"`
	PublicationValue   string                     `json:"publicationValue"`
	ExternalRecipient  string                     `json:"externalRecipient"`
	Subject            string                     `json:"subject"`
	Body               string                     `json:"body"`
	Type               int                        `json:"type"`
	Step               int                        `json:"step"`
	Status             int                        `json:"status"`
	Priority           int                        `json:"priority"`
	Author             model.User                 `json:"author"`
	DocumentSequence   DocumentInProgressResponse `json:"documentSequence"`
	DocumentHistory    *[]DocumentHistory         `json:"documentHistory"`
	DocumentAttachment *[]DocumentAttachment      `json:"documentAttachment"`
	InternalRecipients *[]InternalRecipient       `json:"internalRecipients"`
	DocumentReferences *[]DocumentReference       `json:"documentReferences"`
	CreatedAt          time.Time                  `json:"createdAt"`
	UpdatedAt          time.Time                  `json:"updatedAt"`
	DueDate            *time.Time                 `json:"dueDate"`
	IsApprover         bool                       `json:"isApprover"`
	IsAllowToUpdate    bool                       `json:"isAllowToUpdate"`
	CanRecall          bool                       `json:"canRecall"`
}

type DocumentReference struct {
	Id      string `json:"id"`
	Subject string `json:"subject"`
}

type InternalRecipient struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

type DocumentHistory struct {
	Name       string `json:"name"`
	Title      string `json:"title"`
	IsApproved bool   `json:"isApproved"`
	Reason     string `json:"reason"`
	UpdatedAt  string `json:"updatedAt"`
}

type DocumentAttachment struct {
	Id           string `json:"id"`
	OriginalName string `json:"originalName"`
	FileName     string `json:"fileName"`
	Path         string `json:"path"`
	Size         string `json:"size"`
	Type         string `json:"type"`
}
