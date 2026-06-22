package documentsequence

import (
	"Microservice/helper"
	"Microservice/model"
)

type NumberingGroupRepository interface {
	Create(data model.NumberingGroup) *helper.ErrorModel
	Get(id string, orgID string) (*model.NumberingGroup, *helper.ErrorModel)
	GetAll(orgID string) ([]model.NumberingGroup, *helper.ErrorModel)
	Delete(id string, orgID string) *helper.ErrorModel
}
