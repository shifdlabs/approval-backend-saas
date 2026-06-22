package documentsequence

import (
	"Microservice/helper"
	"Microservice/model"

	"gorm.io/gorm"
)

type DocumentSequenceRepository interface {
	Create(db *gorm.DB, document model.DocumentSequence) *helper.ErrorModel
	Get(id string, orgID string) (*model.DocumentSequence, *helper.ErrorModel)
	GetAll() ([]model.DocumentSequence, *helper.ErrorModel)
	GetSequencesByDocumentId(id string) ([]model.DocumentSequence, *helper.ErrorModel)
	GetAllSequenceByDocumentId(id string) ([]model.DocumentSequence, *helper.ErrorModel)
	GetProgressByAuthorID(authorID string, orgID string) ([]model.DocumentSequence, *helper.ErrorModel)
	GetCurrentApprover(docId string) (*model.DocumentSequence, *helper.ErrorModel)
	Update(document model.Document, sequences []model.DocumentSequence) error
	Delete(id string) *helper.ErrorModel
}
