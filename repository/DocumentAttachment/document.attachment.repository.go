package documentattachment

import (
	"Microservice/helper"
	"Microservice/model"

	"gorm.io/gorm"
)

type DocumentAttachmentRepository interface {
	Create(db *gorm.DB, document model.DocumentAttachment) *helper.ErrorModel
	Get(id string, orgID string) (*model.DocumentAttachment, *helper.ErrorModel)
	GetAll(orgID string) ([]model.DocumentAttachment, *helper.ErrorModel)
	Delete(id string, orgID string) *helper.ErrorModel
}
