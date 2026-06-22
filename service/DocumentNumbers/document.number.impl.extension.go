package documentnumbers

import (
	response "Microservice/data/response/DocumentNumbers"
	"Microservice/model"
)

func (t DocumentNumbersServiceImpl) mapDocumentNumbersToDocumentNumbersResponse(documentHistories []model.DocumentNumbers) []response.DocumentNumbersResponse {
	responseDocuments := make([]response.DocumentNumbersResponse, len(documentHistories))
	for i, documentNumbers := range documentHistories {
		responseDocuments[i] = t.convertDocumentNumbersToDocumentNumbersResponse(documentNumbers)
	}
	return responseDocuments
}

func (t DocumentNumbersServiceImpl) convertDocumentNumbersToDocumentNumbersResponse(documentNumbers model.DocumentNumbers) response.DocumentNumbersResponse {
	responseDocument := response.DocumentNumbersResponse{
		Id:                  documentNumbers.ID,
		DocumentNumber:      documentNumbers.Value,
		NumberingFormatName: documentNumbers.NumberingFormat.Name,
		NumberingFormatId:   documentNumbers.NumberingFormat.ID,
		NumberingGroupName:  documentNumbers.NumberingFormat.Group.Name,
		CreatedAt:           *documentNumbers.CreatedAt,
	}

	return responseDocument
}
