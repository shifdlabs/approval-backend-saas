package appsettings

import (
	request "Microservice/data/request/AppSettings"
	response "Microservice/data/response/AppSettings"
	"Microservice/helper"
	"Microservice/model"
	repository "Microservice/repository/AppSettings"

	"github.com/go-playground/validator/v10"
)

type AppSettingsServiceImpl struct {
	AppSettingsRepository repository.AppSettingsRepository
	Validate              *validator.Validate
}

func NewAppSettingsServiceImpl(
	reportRepository repository.AppSettingsRepository,
	validate *validator.Validate) AppSettingService {
	return &AppSettingsServiceImpl{
		AppSettingsRepository: reportRepository,
		Validate:              validate,
	}
}

func (t AppSettingsServiceImpl) GetAll(orgID string) ([]response.AppSettingResponse, *helper.ErrorModel) {
	result, errFetch := t.AppSettingsRepository.GetAll(orgID)

	if errFetch != nil {
		return nil, errFetch
	} else {
		return t.mapAppSettingsToAppSettingsResponse(result), nil
	}
}

func (t AppSettingsServiceImpl) Update(appSettings request.AppSettingRequest, orgID string) *helper.ErrorModel {
	errStructure := t.Validate.Struct(appSettings)
	if errStructure != nil {
		msg := "Structure Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	var appSettingsModels []model.AppSettings

	for _, value := range appSettings.Properties {
		appSettingsModels = append(appSettingsModels, model.AppSettings{Key: value.Key, Value: value.Value})
	}

	errUpdate := t.AppSettingsRepository.Update(appSettingsModels, orgID)
	if errUpdate != nil {
		return errUpdate
	}

	return nil
}
