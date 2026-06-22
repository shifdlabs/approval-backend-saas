package documentreference

import (
	"Microservice/helper"
	"Microservice/model"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type DocumentReferenceRepository interface {
	Create(db *gorm.DB, data model.DocumentReference) *helper.ErrorModel
	GetAll(documentID uuid.UUID) ([]model.DocumentReference, *helper.ErrorModel)
	Update(newData []string, documentID uuid.UUID) *helper.ErrorModel
	Delete(id uuid.UUID) *helper.ErrorModel
}
