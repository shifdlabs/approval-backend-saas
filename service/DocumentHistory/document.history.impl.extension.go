package documentHistory

import (
	response "Microservice/data/response/DocumentHistory"
	"Microservice/model"
)

func (t DocumentHistoryServiceImpl) mapDocumentHistoryToDocumentHistoryResponse(documentHistories []model.DocumentHistory) []response.DocumentHistoryResponse {
	responseDocuments := make([]response.DocumentHistoryResponse, len(documentHistories))
	for i, documentHistory := range documentHistories {
		responseDocuments[i] = t.convertDocumentHistoryToDocumentHistoryResponse(documentHistory)
	}
	return responseDocuments
}

func (t DocumentHistoryServiceImpl) convertDocumentHistoryToDocumentHistoryResponse(documentHistory model.DocumentHistory) response.DocumentHistoryResponse {
	// Perform necessary conversion logic here, potentially selecting specific fields
	responseDocument := response.DocumentHistoryResponse{
		Id:          documentHistory.ID,
		DocumentID:  &documentHistory.DocumentID,
		UserID:      &documentHistory.UserID,
		Description: documentHistory.Description,
	}

	return responseDocument
}
