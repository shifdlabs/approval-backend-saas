package user

import (
	"Microservice/helper"
	"Microservice/model"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	Db *gorm.DB
}

func NewUserRepositoryImpl(Db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{Db: Db}
}

func (t *UserRepositoryImpl) Create(user model.User) *helper.ErrorModel {
	result := t.Db.Create(&user)

	if result.Error != nil {
		msg := "Create Position Failed"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

func (t *UserRepositoryImpl) Get(id string, hidePassword bool, orgID string) (*model.User, *helper.ErrorModel) {
	var user model.User

	userId, errParse := uuid.Parse(id)
	if errParse != nil {
		msg := "Failed to Parse UUID"
		return nil, helper.ErrorCatcher(errParse, 500, &msg)
	}

	result := t.Db.Preload("Position").Where("organization_id = ?", orgID).First(&user, "id = ?", userId)

	if result.Error != nil {
		msg := "User nout found"
		return nil, helper.ErrorCatcher(result.Error, 404, &msg)
	}

	if hidePassword {
		user.Password = ""
	}

	return &user, nil
}

func (t *UserRepositoryImpl) GetAll(orgID string) ([]model.User, *helper.ErrorModel) {
	var users []model.User
	result := t.Db.Preload("Position").Where("organization_id = ?", orgID).Order("created_at DESC").Find(&users)
	if result.Error != nil {
		msg := "Users not found"
		return nil, helper.ErrorCatcher(result.Error, 404, &msg)
	}

	return users, nil
}

func (t *UserRepositoryImpl) GetAllUserExceptCurrent(userId string, orgID string) ([]model.User, *helper.ErrorModel) {
	var users []model.User
	result := t.Db.Preload("Position").
		Where("organization_id = ?", orgID).
		Where("id != ?", userId).
		Where("role != 99").
		Order("created_at DESC").Find(&users)
	if result.Error != nil {
		msg := "Users not found"
		return nil, helper.ErrorCatcher(result.Error, 404, &msg)
	}

	return users, nil
}

func (t *UserRepositoryImpl) Update(user model.User, orgID string) *helper.ErrorModel {
	var existing model.User
	if err := t.Db.Where("organization_id = ?", orgID).First(&existing, user.ID).Error; err != nil {
		msg := "User not found"
		return helper.ErrorCatcher(err, 404, &msg)
	}

	// We have to add .Select("*") so gorm will not ignoring zero value like 'false', and it can still updating all value
	result := t.Db.Model(&existing).Select("*").Updates(user)
	if result.Error != nil {
		msg := "Failed to Update User Data"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

// Think about how to make it didn't force delete, just set the deleted_At value, but the unique value is not calculated
func (t *UserRepositoryImpl) Delete(id string, orgID string) *helper.ErrorModel {
	userId, errParse := uuid.Parse(id)
	if errParse != nil {
		msg := "Failed to Parse UUID"
		return helper.ErrorCatcher(errParse, 500, &msg)
	}

	result := t.Db.Unscoped().Where("organization_id = ?", orgID).Delete(&model.User{}, userId)

	if result.Error != nil {
		msg := "Failed to Delete User Data"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

func (t *UserRepositoryImpl) MultipleDelete(ids []string, orgID string) *helper.ErrorModel {
	var userIds []uuid.UUID

	for _, id := range ids {
		userId, errParse := uuid.Parse(id)
		if errParse != nil {
			msg := "Failed to Parse UUID"
			return helper.ErrorCatcher(errParse, 500, &msg)
		}
		userIds = append(userIds, userId)
	}

	result := t.Db.Unscoped().Where("id IN ? AND organization_id = ?", userIds, orgID).Delete(&model.User{})

	fmt.Println("Rows affected:", result.RowsAffected)

	if result.Error != nil {
		msg := "Failed to Delete User Data"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

func (t *UserRepositoryImpl) GetByEmail(email string, orgID string) (*model.User, *helper.ErrorModel) {
	var user model.User

	result := t.Db.Preload("Position").Model(model.User{}).Where("organization_id = ?", orgID).First(&user, "email = ?", strings.ToLower(email))

	if result.Error != nil {
		msg := "User nout found"
		return nil, helper.ErrorCatcher(result.Error, 404, &msg)
	}

	return &user, nil
}

// GetUnscoped looks up a user by id with no organization filter. Only for
// pre-authentication flows (see interface doc comment).
func (t *UserRepositoryImpl) GetUnscoped(id string, hidePassword bool) (*model.User, *helper.ErrorModel) {
	var user model.User

	userId, errParse := uuid.Parse(id)
	if errParse != nil {
		msg := "Failed to Parse UUID"
		return nil, helper.ErrorCatcher(errParse, 500, &msg)
	}

	result := t.Db.Preload("Position").First(&user, "id = ?", userId)
	if result.Error != nil {
		msg := "User nout found"
		return nil, helper.ErrorCatcher(result.Error, 404, &msg)
	}

	if hidePassword {
		user.Password = ""
	}

	return &user, nil
}

// GetByEmailUnscoped looks up a user by email with no organization filter.
// Only for pre-authentication flows (see interface doc comment).
func (t *UserRepositoryImpl) GetByEmailUnscoped(email string) (*model.User, *helper.ErrorModel) {
	var user model.User

	result := t.Db.Preload("Position").Model(model.User{}).First(&user, "email = ?", strings.ToLower(email))

	if result.Error != nil {
		msg := "User nout found"
		return nil, helper.ErrorCatcher(result.Error, 404, &msg)
	}

	return &user, nil
}
