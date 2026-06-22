package appsettings

import (
	"Microservice/helper"
	"Microservice/model"
)

type AppSettingsRepository interface {
	GetAll(orgID string) ([]model.AppSettings, *helper.ErrorModel)
	GetByKey(key string, orgID string) (*model.AppSettings, *helper.ErrorModel)
	Update(properties []model.AppSettings, orgID string) *helper.ErrorModel
}
