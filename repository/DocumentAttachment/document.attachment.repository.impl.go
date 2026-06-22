package documentattachment

import (
	"Microservice/helper"
	"Microservice/model"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentAttachmentRepositoryImpl struct {
	Db *gorm.DB
}

func NewDocumentAttachmentRepositoryImpl(Db *gorm.DB) DocumentAttachmentRepository {
	return &DocumentAttachmentRepositoryImpl{Db: Db}
}

func (t *DocumentAttachmentRepositoryImpl) Create(db *gorm.DB, document model.DocumentAttachment) *helper.ErrorModel {
	result := db.Create(&document)

	if result.Error != nil {
		msg := "Failed to create document attachment"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

// document_attachments has no organization_id of its own — it inherits the
// org through its parent documents row (document_id), so scoping joins there.
func (t *DocumentAttachmentRepositoryImpl) Get(id string, orgID string) (*model.DocumentAttachment, *helper.ErrorModel) {
	var documentAttachment model.DocumentAttachment
	documentAttachmentId, err := uuid.Parse(id)
	if err != nil {
		msg := "Failed to parse uuid"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.
		Joins("JOIN documents ON documents.id = document_attachments.document_id").
		Where("documents.organization_id = ?", orgID).
		First(&documentAttachment, "document_attachments.id = ?", documentAttachmentId)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "record not found") {
			return nil, nil
		}

		msg := "Get Document attachment failed"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return &documentAttachment, nil
}

func (t *DocumentAttachmentRepositoryImpl) GetAll(orgID string) ([]model.DocumentAttachment, *helper.ErrorModel) {
	var documentAttachments []model.DocumentAttachment
	result := t.Db.
		Joins("JOIN documents ON documents.id = document_attachments.document_id").
		Where("documents.organization_id = ?", orgID).
		Find(&documentAttachments)
	if result.Error != nil {
		msg := "Failed to get all documents attachments"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return documentAttachments, nil
}

func (t *DocumentAttachmentRepositoryImpl) Delete(id string, orgID string) *helper.ErrorModel {
	documentAttachmentId, err := uuid.Parse(id)
	if err != nil {
		msg := "Failed to parse id"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	// Verify ownership before deleting (standard SQL DELETE has no portable
	// JOIN syntax, so check via SELECT first).
	var existing model.DocumentAttachment
	if err := t.Db.
		Joins("JOIN documents ON documents.id = document_attachments.document_id").
		Where("documents.organization_id = ?", orgID).
		First(&existing, "document_attachments.id = ?", documentAttachmentId).Error; err != nil {
		msg := "Document attachment not found"
		return helper.ErrorCatcher(err, 404, &msg)
	}

	result := t.Db.Unscoped().Delete(&model.DocumentAttachment{}, documentAttachmentId)
	if result.Error != nil {
		msg := "Failed to delete document attachments"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}
