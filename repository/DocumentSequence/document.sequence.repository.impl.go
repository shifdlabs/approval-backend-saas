package documentsequence

import (
	"Microservice/helper"
	"Microservice/model"
	"errors"
	"strings"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DocumentSequenceRepositoryImpl struct {
	Db *gorm.DB
}

func NewDocumentSequenceRepositoryImpl(Db *gorm.DB) DocumentSequenceRepository {
	return &DocumentSequenceRepositoryImpl{Db: Db}
}

func (t *DocumentSequenceRepositoryImpl) Create(db *gorm.DB, document model.DocumentSequence) *helper.ErrorModel {
	result := db.Create(&document)

	if result.Error != nil {
		db.Rollback()
		msg := "Failed to create document sequence"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

func (t *DocumentSequenceRepositoryImpl) Get(id string, orgID string) (*model.DocumentSequence, *helper.ErrorModel) {
	var documentSequence model.DocumentSequence
	documentSequenceId, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse uuid"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.
		Joins("JOIN documents ON documents.id = document_sequences.document_id").
		Where("documents.organization_id = ?", orgID).
		First(&documentSequence, "document_sequences.id = ?", documentSequenceId)

	if result.Error != nil && strings.Contains(result.Error.Error(), "record not found") {
		return nil, nil
	}

	if result.Error != nil {
		msg := "Get Document Sequence Failed"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return &documentSequence, nil
}

func (t *DocumentSequenceRepositoryImpl) GetAll() ([]model.DocumentSequence, *helper.ErrorModel) {
	var documentSequences []model.DocumentSequence
	result := t.Db.Find(&documentSequences)
	if result.Error != nil {
		msg := "Failed to get all documents type"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return documentSequences, nil
}

func (t *DocumentSequenceRepositoryImpl) GetSequencesByDocumentId(id string) ([]model.DocumentSequence, *helper.ErrorModel) {
	var documentSequences []model.DocumentSequence
	documentId, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse id"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	// Get Document Sequence by Document ID where step of sequence >= document.step
	response := t.Db.
		Model(&model.DocumentSequence{}).
		Joins("JOIN documents ON documents.id = document_sequences.document_id").
		Where("document_sequences.document_id = ?", documentId).
		Find(&documentSequences)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			// Return empty Document rather than nil pointer
			return documentSequences, nil
		}

		msg := "Failed to get all documents type"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return documentSequences, nil
}

func (t *DocumentSequenceRepositoryImpl) GetAllSequenceByDocumentId(id string) ([]model.DocumentSequence, *helper.ErrorModel) {
	var documentSequences []model.DocumentSequence
	documentId, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse id"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	// Get Document Sequence by Document ID where step of sequence >= document.step
	response := t.Db.
		Model(&model.DocumentSequence{}).
		Joins("JOIN documents ON documents.id = document_sequences.document_id").
		Where("document_sequences.document_id = ?", documentId).
		Where("document_sequences.step > documents.step").
		Find(&documentSequences)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			// Return empty Document rather than nil pointer
			return documentSequences, nil
		}

		msg := "Failed to get all documents type"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return documentSequences, nil
}

func (t *DocumentSequenceRepositoryImpl) GetCurrentApprover(docId string) (*model.DocumentSequence, *helper.ErrorModel) {
	var documentSequences model.DocumentSequence
	documentId, err := uuid.FromString(docId)
	if err != nil {
		msg := "Failed to parse id"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	// Get Document Sequence by Document ID where step of sequence >= document.step
	response := t.Db.
		Model(&model.DocumentSequence{}).
		Joins("JOIN documents ON documents.id = document_sequences.document_id").
		Where("document_sequences.document_id = ?", documentId).
		Where("document_sequences.step = documents.step").
		Where("documents.status = 1 OR documents.status = 99").
		First(&documentSequences)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			// Return empty Document rather than nil pointer
			return &documentSequences, nil
		}

		msg := "Failed to get all documents type"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return &documentSequences, nil
}

func (t *DocumentSequenceRepositoryImpl) Update(document model.Document, sequences []model.DocumentSequence) error {
	// Use a transaction to ensure atomicity:contentReference[oaicite:11]{index=11}.
	return t.Db.Transaction(func(tx *gorm.DB) error {
		// Deduplicate input by UserID (keeping only the last occurrence).
		seen := make(map[uuid.UUID]bool, len(sequences))
		uniqueCC := make([]model.DocumentSequence, 0, len(sequences))
		for i := len(sequences) - 1; i >= 0; i-- {
			cc := sequences[i]
			if cc.UserID == uuid.Nil {
				continue // skip entries with invalid user UUID
			}
			if !seen[cc.UserID] {
				seen[cc.UserID] = true
				uniqueCC = append(uniqueCC, cc)
			}
		}

		// (Now uniqueCC has at most one entry per UserID, in reverse order of appearance.)

		// Fetch existing user IDs from the DB for this document.
		var existing []model.DocumentSequence
		if err := tx.Model(&model.DocumentSequence{}).
			Where("document_id = ?", document.ID).
			Find(&existing).Error; err != nil {
			return err
		}
		existingSet := make(map[uuid.UUID]bool, len(existing))
		for _, cc := range existing {
			existingSet[cc.UserID] = true
		}

		// Prepare sets of IDs to add and remove.
		var toAddIDs []model.DocumentSequence
		newSet := make(map[uuid.UUID]bool, len(uniqueCC))
		for _, cc := range uniqueCC {
			newSet[cc.UserID] = true
			if !existingSet[cc.UserID] {
				toAddIDs = append(toAddIDs, cc)
			}
		}

		var toRemoveIDs []uuid.UUID
		for _, cc := range existing {
			if !newSet[cc.UserID] {
				toRemoveIDs = append(toRemoveIDs, cc.UserID)
			}
		}

		// Insert new CarbonCopy records in bulk, if any:contentReference[oaicite:12]{index=12}.
		if len(toAddIDs) > 0 {
			newEntries := make([]model.DocumentSequence, 0, len(toAddIDs))
			for _, sequence := range toAddIDs {
				newEntries = append(newEntries, model.DocumentSequence{
					Document:  &document,
					UserID:    sequence.UserID,
					Step:      sequence.Step,
					Signature: sequence.Signature,
				})
			}
			// OnConflict DoNothing prevents duplicate-key errors under race:contentReference[oaicite:13]{index=13}.
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).
				Create(&newEntries).Error; err != nil {
				return err
			}
		}

		// Delete obsolete records in bulk, if any:contentReference[oaicite:14]{index=14}.
		if len(toRemoveIDs) > 0 {
			if err := tx.Unscoped().Where("document_id = ? AND user_id IN ?", document.ID, toRemoveIDs).
				Delete(&model.DocumentSequence{}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (t *DocumentSequenceRepositoryImpl) BulkDelete(ids []string) *helper.ErrorModel {
	var uuids []uuid.UUID

	for _, id := range ids {
		carbonCopyId, errParse := uuid.FromString(id)
		if errParse != nil {
			t.Db.Rollback()
			msg := "Failed to Parse UUID: " + id
			return helper.ErrorCatcher(errParse, 500, &msg)
		}
		uuids = append(uuids, carbonCopyId)
	}

	// Perform a bulk delete using the "IN" condition on the UUID slice
	result := t.Db.Unscoped().Where("user_id IN (?)", uuids).Delete(&model.CarbonCopy{})
	if result.Error != nil {
		t.Db.Rollback()
		msg := "Failed to Delete CarbonCopys Data in bulk"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	t.Db.Commit()

	return nil
}

func (t *DocumentSequenceRepositoryImpl) Delete(id string) *helper.ErrorModel {
	documentSequenceId, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse id"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.Unscoped().Delete(&model.DocumentSequence{}, documentSequenceId)
	if result.Error != nil {
		msg := "Failed to delete document type"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	return nil
}

func (t *DocumentSequenceRepositoryImpl) GetProgressByAuthorID(authorID string, orgID string) ([]model.DocumentSequence, *helper.ErrorModel) {
	var documentSequences []model.DocumentSequence

	// Gunakan Preload untuk memuat relasi User dan Document
	result := t.Db.Preload("Document").Preload("Document.Author").
		Joins("JOIN documents ON documents.id = document_sequences.document_id").
		Where("documents.organization_id = ?", orgID).
		Where("document_sequences.user_id = ? AND documents.status = ?", authorID, 1). // Tambahkan kembali filter status = 1
		Order("document_sequences.step ASC").
		Find(&documentSequences)

	if result.Error != nil {
		msg := "Failed to fetch progress documents for the author"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return documentSequences, nil
}
