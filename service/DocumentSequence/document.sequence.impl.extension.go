package documentSequence

import (
	response "Microservice/data/response/DocumentSequence"
	"Microservice/model"
)

func (t DocumentSequenceServiceImpl) mapDocumentSequenceToDocumentSequenceResponse(documentHistories []model.DocumentSequence) []response.DocumentSequenceResponse {
	responseDocuments := make([]response.DocumentSequenceResponse, len(documentHistories))
	for i, documentSequence := range documentHistories {
		responseDocuments[i] = t.convertDocumentSequenceToDocumentSequenceResponse(documentSequence)
	}
	return responseDocuments
}

func (t DocumentSequenceServiceImpl) convertDocumentSequenceToDocumentSequenceResponse(documentSequence model.DocumentSequence) response.DocumentSequenceResponse {
	// Perform necessary conversion logic here, potentially selecting specific fields
	responseDocument := response.DocumentSequenceResponse{
		Id:         documentSequence.ID,
		Step:       documentSequence.Step,
		UserID:     &documentSequence.UserID,
		DocumentID: documentSequence.DocumentID,
	}

	return responseDocument
}
