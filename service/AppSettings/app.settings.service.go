package appsettings

import (
	request "Microservice/data/request/AppSettings"
	response "Microservice/data/response/AppSettings"
	"Microservice/helper"
)

type AppSettingService interface {
	GetAll(orgID string) ([]response.AppSettingResponse, *helper.ErrorModel)
	Update(appSettings request.AppSettingRequest, orgID string) *helper.ErrorModel
}
