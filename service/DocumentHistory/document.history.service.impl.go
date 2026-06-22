package documentHistory

import (
	document "Microservice/data/response/Document"
	response "Microservice/data/response/DocumentHistory"
	"Microservice/helper"
	"Microservice/model"
	repository "Microservice/repository/DocumentHistory"
	"errors"

	"github.com/go-playground/validator/v10"
)

type DocumentHistoryServiceImpl struct {
	DocumentHistoryRepository repository.DocumentHistoryRepository
	Validate                  *validator.Validate
}

func NewDocumentHistoryServiceImpl(
	documentRepository repository.DocumentHistoryRepository,
	validate *validator.Validate) DocumentHistoryService {
	return &DocumentHistoryServiceImpl{
		DocumentHistoryRepository: documentRepository,
		Validate:                  validate,
	}
}

func (t DocumentHistoryServiceImpl) Get(id string, orgID string) (*response.DocumentHistoryResponse, *helper.ErrorModel) {
	document, fetchError := t.DocumentHistoryRepository.Get(id, orgID)
	if fetchError != nil {
		return nil, fetchError
	}

	if document == nil {
		return nil, nil
	}

	documentResponse := t.convertDocumentHistoryToDocumentHistoryResponse(*document)

	return &documentResponse, fetchError
}

func (t DocumentHistoryServiceImpl) GetAll(orgID string) ([]response.DocumentHistoryResponse, *helper.ErrorModel) {
	result, fetchError := t.DocumentHistoryRepository.GetAll(orgID)
	if fetchError != nil {
		return nil, fetchError
	} else {
		return t.mapDocumentHistoryToDocumentHistoryResponse(result), nil
	}
}

func (t DocumentHistoryServiceImpl) Delete(id string) *helper.ErrorModel {
	errResponse := t.DocumentHistoryRepository.Delete(id)
	if errResponse != nil {
		return errResponse
	}

	return nil
}

func (t DocumentHistoryServiceImpl) FetchHistoriesByUserID(userID string, orgID string) ([]response.DocumentHistoryResponse, *helper.ErrorModel) {
	// Ambil data dari repository
	documentHistories, fetchError := t.DocumentHistoryRepository.GetHistoriesByAuthorID(userID, orgID)
	if fetchError != nil {
		return nil, fetchError
	}

	// Konversi data ke response
	var responses []response.DocumentHistoryResponse
	for _, history := range documentHistories {
		// Validasi UUID
		if len(history.UserID.String()) != 36 || len(history.DocumentID.String()) != 36 {
			msg := "Invalid UUID length"
			return nil, helper.ErrorCatcher(errors.New(msg), 500, &msg)
		}

		// Validasi pointer Author
		var authorResponse model.User
		if history.Document.Author != nil {
			authorResponse = *history.Document.Author
		}

		// Tambahkan data ke responses
		responses = append(responses, response.DocumentHistoryResponse{
			Id:          history.ID,
			DocumentID:  &history.DocumentID,
			UserID:      &history.UserID,
			IsApproved:  history.IsApproved,
			Description: history.Description,
			Document: document.DocumentResponse{
				Id:      &history.Document.ID,
				Subject: history.Document.Subject,
				Status:  history.Document.Status,
				Author:  &authorResponse,
			},
		})
	}

	return responses, nil
}
