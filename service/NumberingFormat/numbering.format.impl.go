package numberingformat

import (
	request "Microservice/data/request/NumberingFormat"
	response "Microservice/data/response/NumberingFormat"
	"Microservice/helper"
)

type NumberingFormatService interface {
	Create(request request.NumberingFormatRequest, orgID string) *helper.ErrorModel
	GetAll(orgID string) ([]response.NumberingFormatResponse, *helper.ErrorModel)
	GetAllWithGrouped(orgID string) ([]response.NumberingFormatByGroupResponse, *helper.ErrorModel)
	Delete(id string, orgID string) *helper.ErrorModel
}
