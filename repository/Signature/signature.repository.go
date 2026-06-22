package signature

import (
	"Microservice/helper"
	"Microservice/model"
)

type SignatureRepository interface {
	Create(signature *model.Signature, orgID string) *helper.ErrorModel
	Update(signature *model.Signature) *helper.ErrorModel
	Delete(id string) *helper.ErrorModel
	GetAll(orgID string) ([]model.Signature, *helper.ErrorModel)
	GetByUserId(userId string, orgID string) (*model.Signature, *helper.ErrorModel)
	GetByUserIds(userIds []string) ([]model.Signature, *helper.ErrorModel)
}
