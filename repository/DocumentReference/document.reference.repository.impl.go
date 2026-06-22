package documentreference

import (
	"Microservice/helper"
	"Microservice/model"
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type DocumentReferenceRepositoryImpl struct {
	Db *gorm.DB
}

func NewDocumentReferenceRepositoryImpl(Db *gorm.DB) DocumentReferenceRepository {
	return &DocumentReferenceRepositoryImpl{Db: Db}
}

func (t *DocumentReferenceRepositoryImpl) Create(db *gorm.DB, data model.DocumentReference) *helper.ErrorModel {
	result := db.Create(&data)

	if result.Error != nil {
		db.Rollback()
		msg := "Failed to create document reference"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

func (t *DocumentReferenceRepositoryImpl) GetAll(documentID uuid.UUID) ([]model.DocumentReference, *helper.ErrorModel) {
	var referencedDocuments []model.DocumentReference
	result := t.Db.Where("document_id = ?", documentID).Find(&referencedDocuments)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Return empty Document rather than nil pointer
			return referencedDocuments, nil
		}

		msg := "Failed to get all documents type"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return referencedDocuments, nil
}

func (t *DocumentReferenceRepositoryImpl) Update(newData []string, documentID uuid.UUID) *helper.ErrorModel {
	msg := "Updating Document Reference Failed"

	// If no data provided → remove all references for this document
	if len(newData) == 0 {
		if err := t.Db.Where("document_id = ?", documentID).Delete(&model.DocumentReference{}).Error; err != nil {
			return helper.ErrorCatcher(err, 500, &msg)
		}
		return nil
	}

	// 1. Get current stored ReferenceIDs
	var referencedDocuments []model.DocumentReference
	if err := t.Db.Where("document_id = ?", documentID).Find(&referencedDocuments).Error; err != nil {
		return helper.ErrorCatcher(err, 500, &msg)
	}

	// Convert existing to map for fast lookup
	existingMap := make(map[string]bool)
	for _, doc := range referencedDocuments {
		existingMap[doc.ReferenceID.String()] = true
	}

	// Convert newData to map for fast lookup
	newMap := make(map[string]bool)
	for _, id := range newData {
		newMap[id] = true
	}

	// 2. Find which to delete (in DB but not in newData)
	var toDelete []uuid.UUID
	for _, doc := range referencedDocuments {
		if !newMap[doc.ReferenceID.String()] {
			toDelete = append(toDelete, doc.ReferenceID)
		}
	}

	if len(toDelete) > 0 {
		if err := t.Db.
			Where("document_id = ? AND reference_id IN ?", documentID, toDelete).
			Delete(&model.DocumentReference{}).Error; err != nil {
			return helper.ErrorCatcher(err, 500, &msg)
		}
	}

	// 3. Find which to insert (in newData but not in DB)
	var toInsert []model.DocumentReference
	for _, id := range newData {
		if !existingMap[id] {
			refID, err := uuid.FromString(id)
			if err != nil {
				return helper.ErrorCatcher(err, 500, &msg)
			}
			toInsert = append(toInsert, model.DocumentReference{
				DocumentID:  documentID,
				ReferenceID: refID,
			})
		}
	}

	if len(toInsert) > 0 {
		if err := t.Db.Create(&toInsert).Error; err != nil {
			return helper.ErrorCatcher(err, 500, &msg)
		}
	}

	return nil
}

func (t *DocumentReferenceRepositoryImpl) Delete(id uuid.UUID) *helper.ErrorModel {
	result := t.Db.Unscoped().Delete(&model.DocumentReference{}, id)
	if result.Error != nil {
		msg := "Failed to delete document type"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}
