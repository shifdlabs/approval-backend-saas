package position

import (
	"Microservice/helper"
	"Microservice/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PositionRepositoryImpl struct {
	Db *gorm.DB
}

func NewPositionRepositoryImpl(Db *gorm.DB) PositionRepository {
	return &PositionRepositoryImpl{Db: Db}
}

func (t *PositionRepositoryImpl) Create(report model.Position) *helper.ErrorModel {
	result := t.Db.Create(&report)

	if result.Error != nil {
		msg := "Create Position Failed"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

func (t *PositionRepositoryImpl) Get(id string, orgID string) (*model.Position, *helper.ErrorModel) {
	var position model.Position

	positionId, errParse := uuid.Parse(id)
	if errParse != nil {
		msg := "Failed to Parse UUID"
		return nil, helper.ErrorCatcher(errParse, 500, &msg)
	}

	result := t.Db.Where("organization_id = ?", orgID).First(&position, "id = ?", positionId)

	if result.Error != nil {
		msg := "Failed to Get Position Data"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return &position, nil
}

func (t *PositionRepositoryImpl) GetAll(orgID string) ([]model.Position, *helper.ErrorModel) {
	var positions []model.Position

	result := t.Db.Where("organization_id = ?", orgID).Find(&positions)
	if result.Error != nil {
		msg := "Failed to Get All Position Data"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return positions, nil
}

func (t *PositionRepositoryImpl) FindByName(name string, orgID string) (*model.Position, *helper.ErrorModel) {
	var position model.Position

	result := t.Db.Where("organization_id = ? AND LOWER(name) = LOWER(?)", orgID, name).First(&position)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil without error if not found
		}
		msg := "Failed to Find Position by Name"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return &position, nil
}

func (t *PositionRepositoryImpl) Update(position model.Position, orgID string) *helper.ErrorModel {
	var existing model.Position

	if err := t.Db.Where("organization_id = ?", orgID).First(&existing, position.ID).Error; err != nil {
		msg := "Position not found"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.Model(&existing).Updates(position)
	if result.Error != nil {
		msg := "Failed to Update Position Data"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

func (t *PositionRepositoryImpl) Delete(id string, orgID string) *helper.ErrorModel {
	positionId, errParse := uuid.Parse(id)
	if errParse != nil {
		msg := "Failed to Parse UUID"
		return helper.ErrorCatcher(errParse, 500, &msg)
	}

	result := t.Db.Unscoped().Where("organization_id = ?", orgID).Delete(&model.Position{}, positionId)
	if result.Error != nil {
		msg := "Failed to Delete Position Data"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}
