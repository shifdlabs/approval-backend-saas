package documentsequence

import (
	"Microservice/helper"
	"Microservice/model"
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type NumberingGroupRepositoryImpl struct {
	Db *gorm.DB
}

func NewNumberingGroupRepositoryImpl(Db *gorm.DB) NumberingGroupRepository {
	return &NumberingGroupRepositoryImpl{Db: Db}
}

func (t *NumberingGroupRepositoryImpl) Create(data model.NumberingGroup) *helper.ErrorModel {
	result := t.Db.Create(&data)

	if result.Error != nil {
		msg := "Failed to create numbering group"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

func (t *NumberingGroupRepositoryImpl) Get(id string, orgID string) (*model.NumberingGroup, *helper.ErrorModel) {
	var numberingGroup model.NumberingGroup
	numberingGroupId, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse uuid"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.Where("organization_id = ? AND id = ?", orgID, numberingGroupId).First(&numberingGroup)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "CarbonCopys not found"
		return nil, helper.ErrorCatcher(result.Error, 404, &msg)
	}

	return &numberingGroup, nil
}

func (t *NumberingGroupRepositoryImpl) GetAll(orgID string) ([]model.NumberingGroup, *helper.ErrorModel) {
	var numberingGroups []model.NumberingGroup
	result := t.Db.Where("organization_id = ? AND deleted_at IS NULL", orgID).Find(&numberingGroups)
	if result.Error != nil {
		msg := "Failed to get all numbering groups"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return numberingGroups, nil
}

func (t *NumberingGroupRepositoryImpl) Delete(id string, orgID string) *helper.ErrorModel {
	numberingGroupId, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse id"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.Where("organization_id = ?", orgID).Delete(&model.NumberingGroup{}, numberingGroupId)
	if result.Error != nil {
		msg := "Failed to delete numbering group"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	return nil
}
