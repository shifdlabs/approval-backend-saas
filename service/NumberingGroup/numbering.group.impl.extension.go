package numberinggroup

import (
	response "Microservice/data/response/NumberingGroup"
	"Microservice/model"
)

func (t NumberingGroupServiceImpl) mapNumberingGroupToNumberingGroupResponse(documentHistories []model.NumberingGroup) []response.NumberingGroupResponse {
	responseDocuments := make([]response.NumberingGroupResponse, len(documentHistories))
	for i, numberingGroup := range documentHistories {
		responseDocuments[i] = t.convertNumberingGroupToNumberingGroupResponse(numberingGroup)
	}
	return responseDocuments
}

func (t NumberingGroupServiceImpl) convertNumberingGroupToNumberingGroupResponse(numberingGroup model.NumberingGroup) response.NumberingGroupResponse {
	// Perform necessary conversion logic here, potentially selecting specific fields
	responseDocument := response.NumberingGroupResponse{
		Id:          numberingGroup.ID,
		Name:        numberingGroup.Name,
		Description: numberingGroup.Description,
		TotalItem:   0, // Please fix this get the total from the numbering format table
	}

	return responseDocument
}
