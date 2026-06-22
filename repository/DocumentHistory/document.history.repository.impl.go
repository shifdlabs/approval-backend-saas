package documenthistory

import (
	"Microservice/helper"
	"Microservice/model"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentHistoryRepositoryImpl struct {
	Db *gorm.DB
}

func NewDocumentHistoryRepositoryImpl(Db *gorm.DB) DocumentHistoryRepository {
	return &DocumentHistoryRepositoryImpl{Db: Db}
}

func (t *DocumentHistoryRepositoryImpl) Create(document model.DocumentHistory) *helper.ErrorModel {
	result := t.Db.Create(&document)

	if result.Error != nil {
		msg := "Failed to create document type"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

// document_histories has no organization_id of its own — it inherits the org
// through its parent documents row (document_id), so scoping joins there.
func (t *DocumentHistoryRepositoryImpl) Get(id string, orgID string) (*model.DocumentHistory, *helper.ErrorModel) {
	var documentHistory model.DocumentHistory

	documentHistoryId, err := uuid.Parse(id)
	if err != nil {
		msg := "Get Document History Failed"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.
		Joins("JOIN documents ON documents.id = document_histories.document_id").
		Where("documents.organization_id = ?", orgID).
		First(&documentHistory, "document_histories.id = ?", documentHistoryId)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "record not found") {
			return nil, nil
		}

		msg := "Get Document History Failed"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return &documentHistory, nil
}

func (t *DocumentHistoryRepositoryImpl) GetAll(orgID string) ([]model.DocumentHistory, *helper.ErrorModel) {
	var documentHistorys []model.DocumentHistory
	result := t.Db.
		Joins("JOIN documents ON documents.id = document_histories.document_id").
		Where("documents.organization_id = ?", orgID).
		Find(&documentHistorys)
	if result.Error != nil {
		msg := "Failed to get all documents type"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return documentHistorys, nil
}

func (t *DocumentHistoryRepositoryImpl) GetAllHistoryByDocumentId(id string) ([]model.DocumentHistory, *helper.ErrorModel) {
	var documentHistorys []model.DocumentHistory

	documentId, err := uuid.Parse(id)
	if err != nil {
		msg := "Failed to parse id"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.Where("document_id = ?", documentId).Find(&documentHistorys)

	if result.Error != nil {
		msg := "Failed to get all documents type"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return documentHistorys, nil
}

func (t *DocumentHistoryRepositoryImpl) GetLastRejection(id string) (*model.DocumentHistory, *helper.ErrorModel) {
	var documentHistory model.DocumentHistory

	documentId, err := uuid.Parse(id)
	if err != nil {
		msg := "Failed to parse id"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	response := t.Db.
		Where("document_id = ?", documentId).
		Where("is_approved = ?", false).
		Order("created_at desc").
		First(&documentHistory)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "Failed to get all documents type"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return &documentHistory, nil
}

func (t *DocumentHistoryRepositoryImpl) GetLastApprover(id string) (*model.DocumentHistory, *helper.ErrorModel) {
	var documentHistory model.DocumentHistory

	documentId, err := uuid.Parse(id)
	if err != nil {
		msg := "Failed to parse id"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	response := t.Db.
		Where("document_id = ?", documentId).
		Order("created_at desc").
		First(&documentHistory)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "Failed to get all documents type"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return &documentHistory, nil
}

func (t *DocumentHistoryRepositoryImpl) Delete(id string) *helper.ErrorModel {
	documentHistoryId, err := uuid.Parse(id)
	if err != nil {
		msg := "Failed to parse id"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.Unscoped().Delete(&model.DocumentHistory{}, documentHistoryId)

	if result.Error != nil {
		msg := "Failed to delete document type"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	return nil
}

func (t *DocumentHistoryRepositoryImpl) GetHistoriesByAuthorID(authorID string, orgID string) ([]model.DocumentHistory, *helper.ErrorModel) {
	var documentHistories []model.DocumentHistory

	// Query untuk mengambil data berdasarkan AuthorID dari Document
	result := t.Db.Preload("Document").Preload("Document.Author").
		Joins("JOIN documents ON documents.id = document_histories.document_id").
		Where("documents.organization_id = ?", orgID).
		Where("documents.status = ? AND documents.author_id = ?", 99, authorID).
		Order("document_histories.created_at ASC").
		Find(&documentHistories)

	if result.Error != nil {
		msg := "Failed to fetch document histories by author ID"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return documentHistories, nil
}
