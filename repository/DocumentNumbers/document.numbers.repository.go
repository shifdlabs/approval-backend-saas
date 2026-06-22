package documentnumbers

import (
	"Microservice/helper"
	"Microservice/model"

	uuid "github.com/satori/go.uuid"
)

type DocumentNumbersRepository interface {
	Create(data model.DocumentNumbers) *helper.ErrorModel
	Get(id string, orgID string) (*model.DocumentNumbers, *helper.ErrorModel)
	GetByDocumentID(id uuid.UUID, orgID string) (*model.DocumentNumbers, *helper.ErrorModel)
	GetAll(orgID string) ([]model.DocumentNumbers, *helper.ErrorModel)
	GetAllByUserID(userId string, orgID string) ([]model.DocumentNumbers, *helper.ErrorModel)
	GetTotal(formatId string, groupId *string, orgID string) (*int64, *helper.ErrorModel)
	GetCancelled(formatId string, groupId *string, orgID string) (*model.DocumentNumbers, *helper.ErrorModel)
	Update(data model.DocumentNumbers, orgID string) *helper.ErrorModel
	Delete(id string, orgID string) *helper.ErrorModel
}
