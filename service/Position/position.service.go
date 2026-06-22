package position

import (
	request "Microservice/data/request/Position"
	response "Microservice/data/response/Position"
	"Microservice/helper"
)

type PositionService interface {
	Create(position request.CreatePositionRequest, orgID string) *helper.ErrorModel
	Get(id string, orgID string) (*response.PositionResponse, *helper.ErrorModel)
	GetAll(orgID string) ([]response.PositionResponse, *helper.ErrorModel)
	Update(position request.UpdatePositionRequest, orgID string) *helper.ErrorModel
	Delete(id string, orgID string) *helper.ErrorModel
}
