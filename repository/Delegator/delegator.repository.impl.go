package delegator

import (
	"Microservice/helper"
	"Microservice/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DelegatorRepositoryImpl struct {
	Db *gorm.DB
}

func NewDelegatorRepositoryImpl(db *gorm.DB) DelegatorRepository {
	return &DelegatorRepositoryImpl{Db: db}
}

// delegators has no organization_id of its own — it inherits the org through
// its owner/delegate users, so scoping joins to the users table.

func (r *DelegatorRepositoryImpl) Create(delegator model.Delegator, orgID string) *helper.ErrorModel {
	// Verify both owner and delegate belong to this org before creating —
	// otherwise a delegation could be assigned across organizations.
	var count int64
	if err := r.Db.Model(&model.User{}).
		Where("organization_id = ?", orgID).
		Where("id IN ?", []string{delegator.OwnerID.String(), delegator.DelegateID.String()}).
		Count(&count).Error; err != nil {
		msg := "Failed to verify owner/delegate"
		return helper.ErrorCatcher(err, 500, &msg)
	}
	if count != 2 {
		msg := "Owner or delegate does not belong to this organization"
		return helper.ErrorCatcher(gorm.ErrRecordNotFound, 400, &msg)
	}

	result := r.Db.Create(&delegator)
	if result.Error != nil {
		msg := "Failed to create delegation"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return nil
}

func (r *DelegatorRepositoryImpl) GetAllByOwnerID(ownerID string, orgID string) ([]model.Delegator, *helper.ErrorModel) {
	var delegators []model.Delegator
	result := r.Db.Preload("Owner").Preload("Delegate").
		Joins("JOIN users ON users.id = delegators.owner_id").
		Where("users.organization_id = ?", orgID).
		Where("delegators.owner_id = ?", ownerID).
		Order("delegators.created_at DESC").
		Find(&delegators)
	if result.Error != nil {
		msg := "Failed to get delegations"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return delegators, nil
}

func (r *DelegatorRepositoryImpl) GetByID(id string, orgID string) (*model.Delegator, *helper.ErrorModel) {
	var delegator model.Delegator
	delegatorID, errParse := uuid.Parse(id)
	if errParse != nil {
		msg := "Failed to parse UUID"
		return nil, helper.ErrorCatcher(errParse, 400, &msg)
	}
	result := r.Db.Preload("Owner").Preload("Delegate").
		Joins("JOIN users ON users.id = delegators.owner_id").
		Where("users.organization_id = ?", orgID).
		First(&delegator, "delegators.id = ?", delegatorID)
	if result.Error != nil {
		msg := "Delegation not found"
		return nil, helper.ErrorCatcher(result.Error, 404, &msg)
	}
	return &delegator, nil
}

func (r *DelegatorRepositoryImpl) Update(delegator model.Delegator, orgID string) *helper.ErrorModel {
	if _, err := r.GetByID(delegator.ID.String(), orgID); err != nil {
		return err
	}

	result := r.Db.Model(&delegator).Updates(map[string]interface{}{
		"delegate_id": delegator.DelegateID,
		"start_date":  delegator.StartDate,
		"end_date":    delegator.EndDate,
	})
	if result.Error != nil {
		msg := "Failed to update delegation"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return nil
}

func (r *DelegatorRepositoryImpl) Delete(id string, orgID string) *helper.ErrorModel {
	if _, err := r.GetByID(id, orgID); err != nil {
		return err
	}

	delegatorID, errParse := uuid.Parse(id)
	if errParse != nil {
		msg := "Failed to parse UUID"
		return helper.ErrorCatcher(errParse, 400, &msg)
	}
	result := r.Db.Unscoped().Delete(&model.Delegator{}, delegatorID)
	if result.Error != nil {
		msg := "Failed to delete delegation"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return nil
}

func (r *DelegatorRepositoryImpl) GetActiveDelegationByOwnerID(ownerID string, date time.Time, orgID string) (*model.Delegator, *helper.ErrorModel) {
	var delegator model.Delegator
	result := r.Db.Preload("Delegate").
		Joins("JOIN users ON users.id = delegators.owner_id").
		Where("users.organization_id = ?", orgID).
		Where("delegators.owner_id = ? AND delegators.start_date <= ? AND delegators.end_date >= ?", ownerID, date, date).
		First(&delegator)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		msg := "Failed to get active delegation"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return &delegator, nil
}

func (r *DelegatorRepositoryImpl) GetOwnerIDsByDelegateID(delegateID string, date time.Time, orgID string) ([]string, *helper.ErrorModel) {
	var ownerIDs []string
	result := r.Db.Model(&model.Delegator{}).
		Joins("JOIN users ON users.id = delegators.owner_id").
		Select("delegators.owner_id").
		Where("users.organization_id = ?", orgID).
		Where("delegators.delegate_id = ? AND delegators.start_date <= ? AND delegators.end_date >= ?", delegateID, date, date).
		Pluck("delegators.owner_id", &ownerIDs)
	if result.Error != nil {
		msg := "Failed to get owners by delegate"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return ownerIDs, nil
}

func (r *DelegatorRepositoryImpl) HasOverlappingDelegation(ownerID string, startDate time.Time, endDate time.Time, excludeID *string, orgID string) (bool, *helper.ErrorModel) {
	query := r.Db.Model(&model.Delegator{}).
		Joins("JOIN users ON users.id = delegators.owner_id").
		Where("users.organization_id = ?", orgID).
		Where("delegators.owner_id = ? AND delegators.start_date <= ? AND delegators.end_date >= ?", ownerID, endDate, startDate)

	if excludeID != nil {
		query = query.Where("delegators.id != ?", *excludeID)
	}

	var count int64
	result := query.Count(&count)
	if result.Error != nil {
		msg := "Failed to check overlapping delegation"
		return false, helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return count > 0, nil
}
