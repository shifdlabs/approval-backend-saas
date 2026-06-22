package numberingformat

import (
	response "Microservice/data/response/NumberingFormat"
	"Microservice/model"
)

func (t NumberingFormatServiceImpl) mapNumberingFormatToNumberingFormatResponse(documentHistories []model.NumberingFormat) []response.NumberingFormatResponse {
	responseDocuments := make([]response.NumberingFormatResponse, len(documentHistories))
	for i, numberingFormat := range documentHistories {
		responseDocuments[i] = t.convertNumberingFormatToNumberingFormatResponse(numberingFormat)
	}
	return responseDocuments
}

func (t NumberingFormatServiceImpl) convertNumberingFormatToNumberingFormatResponse(numberingFormat model.NumberingFormat) response.NumberingFormatResponse {
	// Perform necessary conversion logic here, potentially selecting specific fields
	responseDocument := response.NumberingFormatResponse{
		Id:        numberingFormat.ID,
		Name:      numberingFormat.Name,
		Format:    numberingFormat.Format,
		Separator: numberingFormat.Separator,
	}

	return responseDocument
}
