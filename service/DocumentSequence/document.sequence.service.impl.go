package documentSequence

import (
	documentResponse "Microservice/data/response/Document" // Untuk DocumentResponse
	response "Microservice/data/response/DocumentSequence" // Untuk DocumentSequenceResponse
	responseUser "Microservice/data/response/User"         // Untuk UserResponse
	"Microservice/helper"
	repository "Microservice/repository/DocumentSequence" // Untuk DocumentSequenceRepository
	"errors"

	"github.com/go-playground/validator/v10"
)

type DocumentSequenceServiceImpl struct {
	DocumentSequenceRepository repository.DocumentSequenceRepository
	Validate                   *validator.Validate
}

func NewDocumentSequenceServiceImpl(
	documentRepository repository.DocumentSequenceRepository,
	validate *validator.Validate) DocumentSequenceService {
	return &DocumentSequenceServiceImpl{
		DocumentSequenceRepository: documentRepository,
		Validate:                   validate,
	}
}

func (t DocumentSequenceServiceImpl) Get(id string, orgID string) (*response.DocumentSequenceResponse, *helper.ErrorModel) {
	document, fetchError := t.DocumentSequenceRepository.Get(id, orgID)
	if fetchError != nil {
		return nil, fetchError
	}

	if document == nil {
		return nil, nil
	}

	documentResponse := t.convertDocumentSequenceToDocumentSequenceResponse(*document)

	return &documentResponse, fetchError
}

func (t DocumentSequenceServiceImpl) GetAll() ([]response.DocumentSequenceResponse, *helper.ErrorModel) {
	result, fetchError := t.DocumentSequenceRepository.GetAll()
	if fetchError != nil {
		return nil, fetchError
	} else {
		return t.mapDocumentSequenceToDocumentSequenceResponse(result), nil
	}
}

func (t DocumentSequenceServiceImpl) Delete(id string) *helper.ErrorModel {
	errResponse := t.DocumentSequenceRepository.Delete(id)
	if errResponse != nil {
		return errResponse
	}

	return nil
}

func (t DocumentSequenceServiceImpl) GetProgressByAuthorID(authorID string, orgID string) ([]response.DocumentSequenceResponse, *helper.ErrorModel) {
	// Ambil data dari repository
	documentSequences, fetchError := t.DocumentSequenceRepository.GetProgressByAuthorID(authorID, orgID)
	if fetchError != nil {
		return nil, fetchError
	}

	// Konversi data ke response
	documentSequenceResponses := make([]response.DocumentSequenceResponse, len(documentSequences))
	for i, documentSequence := range documentSequences {
		// Validasi DocumentID
		if documentSequence.DocumentID == nil || len(documentSequence.DocumentID.String()) != 36 {
			msg := "Invalid UUID length for DocumentID"
			return nil, helper.ErrorCatcher(errors.New(msg), 500, &msg)
		}

		// Validasi Document
		if documentSequence.Document == nil {
			msg := "Document is nil"
			return nil, helper.ErrorCatcher(errors.New(msg), 500, &msg)
		}

		// Validasi Author
		var authorResponse responseUser.UserResponse
		if documentSequence.Document.Author != nil {
			authorResponse = responseUser.UserResponse{
				ID:        documentSequence.Document.Author.ID,
				FirstName: documentSequence.Document.Author.FirstName,
				LastName:  documentSequence.Document.Author.LastName,
				Email:     documentSequence.Document.Author.Email,
			}
		}

		// Tambahkan data ke response
		documentSequenceResponses[i] = response.DocumentSequenceResponse{
			Id:         documentSequence.ID,
			Step:       documentSequence.Step,
			UserID:     &documentSequence.UserID,
			DocumentID: documentSequence.DocumentID,
			User:       authorResponse,
			Document: documentResponse.DocumentResponse{
				Id:      &documentSequence.Document.ID,
				Subject: documentSequence.Document.Subject,
				Status:  documentSequence.Document.Status,
			},
		}
	}

	return documentSequenceResponses, nil
}
