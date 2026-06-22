package numberingformat

import (
	"Microservice/helper"
	"Microservice/model"
)

type NumberingFormatRepository interface {
	Create(data model.NumberingFormat) *helper.ErrorModel
	Get(id string, orgID string) (*model.NumberingFormat, *helper.ErrorModel)
	GetAll(orgID string) ([]model.NumberingFormat, *helper.ErrorModel)
	Delete(id string, orgID string) *helper.ErrorModel
}
