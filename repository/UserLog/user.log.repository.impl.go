package userlog

import (
	"Microservice/helper"
	"Microservice/model"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type UserLogRepositoryImpl struct {
	Db *gorm.DB
}

func NewUserLogRepositoryImpl(Db *gorm.DB) UserLogRepository {
	return &UserLogRepositoryImpl{Db: Db}
}

func (t *UserLogRepositoryImpl) Create(document model.UserLog, orgID string) {
	if orgUUID, err := uuid.FromString(orgID); err == nil {
		document.OrganizationID = &orgUUID
	}

	result := t.Db.Create(&document)

	if result.Error != nil {
		msg := "Failed to create user log"
		helper.ErrorLog(result.Error, 500, &msg)
	}
}

func (t *UserLogRepositoryImpl) GetAll(orgID string) ([]UserLogWithName, *helper.ErrorModel) {
	var rows []UserLogWithName
	result := t.Db.Table("user_logs").
		Select("user_logs.*, COALESCE(u.first_name || ' ' || u.last_name, 'Unknown') AS user_name").
		Joins("LEFT JOIN users u ON u.id = user_logs.user_id AND u.deleted_at IS NULL").
		Where("user_logs.organization_id = ?", orgID).
		Where("user_logs.deleted_at IS NULL").
		Order("user_logs.log_date DESC").
		Scan(&rows)
	if result.Error != nil {
		msg := "Failed to get all user logs"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return rows, nil
}
