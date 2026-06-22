package numberinggroup

import (
	request "Microservice/data/request/NumberingGroup"
	response "Microservice/data/response/NumberingGroup"
	"Microservice/helper"
)

type NumberingGroupService interface {
	Create(request request.NumberingGroupRequest, orgID string) *helper.ErrorModel
	Get(id string, orgID string) (*response.NumberingGroupResponse, *helper.ErrorModel)
	GetAll(orgID string) ([]response.NumberingGroupResponse, *helper.ErrorModel)
	Delete(id string, orgID string) *helper.ErrorModel
}
