package user

import (
	"Microservice/helper"
	"Microservice/model"
)

type UserRepository interface {
	Create(user model.User) *helper.ErrorModel
	Get(id string, hidePassword bool, orgID string) (*model.User, *helper.ErrorModel)
	GetAll(orgID string) ([]model.User, *helper.ErrorModel)
	GetAllUserExceptCurrent(userId string, orgID string) ([]model.User, *helper.ErrorModel)
	GetByEmail(email string, orgID string) (*model.User, *helper.ErrorModel)
	Update(user model.User, orgID string) *helper.ErrorModel
	Delete(id string, orgID string) *helper.ErrorModel
	MultipleDelete(ids []string, orgID string) *helper.ErrorModel

	// GetUnscoped and GetByEmailUnscoped intentionally bypass org filtering.
	// They exist only for pre-authentication flows (password reset, the
	// now-unrouted local login/refresh) where no org_id is known yet — the
	// caller has no JWT, so there is nothing to scope by. Do not use these
	// from any handler that already has an authenticated org_id in context.
	GetUnscoped(id string, hidePassword bool) (*model.User, *helper.ErrorModel)
	GetByEmailUnscoped(email string) (*model.User, *helper.ErrorModel)
}
