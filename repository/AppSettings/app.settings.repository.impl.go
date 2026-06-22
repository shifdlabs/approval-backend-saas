package appsettings

import (
	"Microservice/helper"
	"Microservice/model"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type AppSettingsRepositoryImpl struct {
	Db *gorm.DB
}

func NewAppSettingsRepositoryImpl(Db *gorm.DB) AppSettingsRepository {
	return &AppSettingsRepositoryImpl{Db: Db}
}

func (t *AppSettingsRepositoryImpl) Create(report model.AppSettings) *helper.ErrorModel {
	result := t.Db.Create(&report)

	if result.Error != nil {
		msg := "Create AppSettings Failed"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

func (t *AppSettingsRepositoryImpl) GetAll(orgID string) ([]model.AppSettings, *helper.ErrorModel) {
	var appSettingss []model.AppSettings

	result := t.Db.Where("organization_id = ?", orgID).Find(&appSettingss)
	if result.Error != nil {
		msg := "Failed to Get All AppSettings Data"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return appSettingss, nil
}

func (t *AppSettingsRepositoryImpl) GetByKey(key string, orgID string) (*model.AppSettings, *helper.ErrorModel) {
	var setting model.AppSettings
	result := t.Db.Where("organization_id = ? AND key = ?", orgID, key).First(&setting)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		msg := "Failed to get app setting by key"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return &setting, nil
}

func (t *AppSettingsRepositoryImpl) Update(appSettings []model.AppSettings, orgID string) *helper.ErrorModel {
	orgUUID, errParse := uuid.FromString(orgID)
	if errParse != nil {
		msg := "Invalid Organization ID"
		return helper.ErrorCatcher(errParse, 500, &msg)
	}

	trx := t.Db.Begin()
	trx.Begin()

	for _, value := range appSettings {
		var existing model.AppSettings
		err := t.Db.Where("organization_id = ? AND key = ?", orgID, value.Key).First(&existing).Error
		if err != nil {
			value.OrganizationID = &orgUUID
			if errCreate := t.Db.Create(&value).Error; errCreate != nil {
				msg := "Failed to Get All AppSettings Data"
				return helper.ErrorCatcher(errCreate, 500, &msg)
			}
		} else {
			if errUpdate := t.Db.Model(&existing).Updates(value).Error; err != nil {
				msg := "Failed to Get All AppSettings Data"
				return helper.ErrorCatcher(errUpdate, 500, &msg)
			}
		}
	}

	trx.Commit()
	return nil
}
