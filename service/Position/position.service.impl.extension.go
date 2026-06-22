package position

import (
	response "Microservice/data/response/Position"
	"Microservice/model"
)

func (t PositionServiceImpl) mapPositionToPositionResponse(positions []model.Position) []response.PositionResponse {
	responseReports := make([]response.PositionResponse, len(positions))
	for i, position := range positions {
		responseReports[i] = t.convertPositionToPositionResponse(position)
	}
	return responseReports
}

func (t PositionServiceImpl) convertPositionToPositionResponse(position model.Position) response.PositionResponse {
	// Perform necessary conversion logic here, potentially selecting specific fields
	responseReport := response.PositionResponse{
		Id:   position.ID,
		Name: position.Name,
	}

	return responseReport
}
