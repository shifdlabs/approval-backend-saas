package position

import (
	"Microservice/helper"
	"Microservice/model"
)

type PositionRepository interface {
	Create(report model.Position) *helper.ErrorModel
	Get(id string, orgID string) (*model.Position, *helper.ErrorModel)
	GetAll(orgID string) ([]model.Position, *helper.ErrorModel)
	FindByName(name string, orgID string) (*model.Position, *helper.ErrorModel)
	Update(position model.Position, orgID string) *helper.ErrorModel
	Delete(id string, orgID string) *helper.ErrorModel
}
