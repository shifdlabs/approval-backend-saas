package numberingformat

import (
	"Microservice/helper"
	"Microservice/model"
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type NumberingFormatRepositoryImpl struct {
	Db *gorm.DB
}

func NewNumberingFormatRepositoryImpl(Db *gorm.DB) NumberingFormatRepository {
	return &NumberingFormatRepositoryImpl{Db: Db}
}

func (t *NumberingFormatRepositoryImpl) Create(data model.NumberingFormat) *helper.ErrorModel {
	result := t.Db.Create(&data)

	if result.Error != nil {
		msg := "Failed to create numbering format"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

// numbering_formats has no organization_id of its own — it inherits the org
// through its parent numbering_groups row (group_id), so scoping joins there.
func (t *NumberingFormatRepositoryImpl) Get(id string, orgID string) (*model.NumberingFormat, *helper.ErrorModel) {
	// Return nil if ID is empty
	if id == "" {
		return nil, nil
	}

	var numberingFormat model.NumberingFormat
	numberingFormatId, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse uuid"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.
		Joins("JOIN numbering_groups ON numbering_groups.id = numbering_formats.group_id").
		Where("numbering_groups.organization_id = ?", orgID).
		Where("numbering_formats.id = ?", numberingFormatId).
		First(&numberingFormat)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "CarbonCopys not found"
		return nil, helper.ErrorCatcher(result.Error, 404, &msg)
	}

	return &numberingFormat, nil
}

func (t *NumberingFormatRepositoryImpl) GetAll(orgID string) ([]model.NumberingFormat, *helper.ErrorModel) {
	var numberingFormats []model.NumberingFormat
	result := t.Db.Preload("Group").
		Joins("JOIN numbering_groups ON numbering_groups.id = numbering_formats.group_id").
		Where("numbering_groups.organization_id = ?", orgID).
		Where("numbering_formats.deleted_at IS NULL").
		Find(&numberingFormats)
	if result.Error != nil {
		msg := "Failed to get all numbering formats"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return numberingFormats, nil
}

func (t *NumberingFormatRepositoryImpl) Delete(id string, orgID string) *helper.ErrorModel {
	numberingFormatId, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse id"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	// Verify the row belongs to this org before deleting (standard SQL DELETE
	// has no portable JOIN syntax, so we check ownership with a SELECT first).
	var existing model.NumberingFormat
	if err := t.Db.
		Joins("JOIN numbering_groups ON numbering_groups.id = numbering_formats.group_id").
		Where("numbering_groups.organization_id = ?", orgID).
		Where("numbering_formats.id = ?", numberingFormatId).
		First(&existing).Error; err != nil {
		msg := "Numbering format not found"
		return helper.ErrorCatcher(err, 404, &msg)
	}

	result := t.Db.Delete(&model.NumberingFormat{}, numberingFormatId)
	if result.Error != nil {
		msg := "Failed to delete numbering format"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}
