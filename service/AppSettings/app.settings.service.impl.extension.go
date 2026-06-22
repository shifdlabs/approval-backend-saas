package appsettings

import (
	response "Microservice/data/response/AppSettings"
	"Microservice/model"
)

func (t AppSettingsServiceImpl) mapAppSettingsToAppSettingsResponse(appSettingss []model.AppSettings) []response.AppSettingResponse {
	responseReports := make([]response.AppSettingResponse, len(appSettingss))
	for i, appSettings := range appSettingss {
		responseReports[i] = t.convertAppSettingsToAppSettingsResponse(appSettings)
	}
	return responseReports
}

func (t AppSettingsServiceImpl) convertAppSettingsToAppSettingsResponse(appSettings model.AppSettings) response.AppSettingResponse {
	// Perform necessary conversion logic here, potentially selecting specific fields
	responseReport := response.AppSettingResponse{
		Key:   appSettings.Key,
		Value: appSettings.Value,
	}

	return responseReport
}
