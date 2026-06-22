package failedloginattempt

import (
	"Microservice/helper"
	"Microservice/model"

	"gorm.io/gorm"
)

type FailedLoginAttemptRepositoryImpl struct {
	Db *gorm.DB
}

func NewFailedLoginAttemptRepositoryImpl(db *gorm.DB) FailedLoginAttemptRepository {
	return &FailedLoginAttemptRepositoryImpl{Db: db}
}

func (f *FailedLoginAttemptRepositoryImpl) Create(attempt model.FailedLoginAttempt) *helper.ErrorModel {
	result := f.Db.Create(&attempt)
	if result.Error != nil {
		msg := "Failed to create failed login attempt record"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return nil
}

func (f *FailedLoginAttemptRepositoryImpl) CountByUserId(userId string) (int64, *helper.ErrorModel) {
	var count int64
	result := f.Db.Model(&model.FailedLoginAttempt{}).Where("user_id = ?", userId).Count(&count)
	if result.Error != nil {
		msg := "Failed to count failed login attempts"
		return 0, helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return count, nil
}

func (f *FailedLoginAttemptRepositoryImpl) DeleteByUserId(userId string) *helper.ErrorModel {
	// Use Unscoped() to perform hard delete instead of soft delete
	result := f.Db.Unscoped().Where("user_id = ?", userId).Delete(&model.FailedLoginAttempt{})
	if result.Error != nil {
		msg := "Failed to delete failed login attempts"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return nil
}
