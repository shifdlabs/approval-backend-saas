package documentSequence

import (
	documentResponse "Microservice/data/response/Document"
	response "Microservice/data/response/User"

	uuid "github.com/satori/go.uuid"
)

type DocumentSequenceResponse struct {
	Id         *uuid.UUID                        `json:"id"`
	Step       int                               `json:"step"`
	UserID     *uuid.UUID                        `json:"userId"`
	DocumentID *uuid.UUID                        `json:"documentId"`
	User       response.UserResponse             `json:"user"`
	Document   documentResponse.DocumentResponse `json:"document"`
}
