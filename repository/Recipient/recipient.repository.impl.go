package recipient

import (
	"Microservice/helper"
	"Microservice/model"
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RecipientRepositoryImpl struct {
	Db *gorm.DB
}

func NewRecipientRepositoryImpl(Db *gorm.DB) RecipientRepository {
	return &RecipientRepositoryImpl{Db: Db}
}

func (t *RecipientRepositoryImpl) Create(db gorm.DB, recipients []model.Recipient) *helper.ErrorModel {
	for _, recipient := range recipients {
		if err := db.Create(&recipient).Error; err != nil {
			db.Rollback()
			msg := "Create Recipient Failed"
			return helper.ErrorCatcher(err, 500, &msg)
		}
	}
	return nil
}

func (t *RecipientRepositoryImpl) GetRecipientsByDocId(id string) ([]model.Recipient, *helper.ErrorModel) {
	var recipients []model.Recipient
	response := t.Db.Where("document_id = ?", id).Find(&recipients)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "Recipients not found"
		return nil, helper.ErrorCatcher(response.Error, 404, &msg)
	}

	return recipients, nil
}

func (t *RecipientRepositoryImpl) GetAll() ([]model.Recipient, *helper.ErrorModel) {
	var recipients []model.Recipient
	result := t.Db.Find(&recipients)
	if result.Error != nil {
		msg := "Recipients not found"
		return nil, helper.ErrorCatcher(result.Error, 404, &msg)
	}

	return recipients, nil
}

func (t *RecipientRepositoryImpl) Update(document model.Document, recipients []model.Recipient) error {
	// Use a transaction to ensure atomicity:contentReference[oaicite:11]{index=11}.
	return t.Db.Transaction(func(tx *gorm.DB) error {
		// Deduplicate input by UserID (keeping only the last occurrence).
		seen := make(map[uuid.UUID]bool, len(recipients))
		uniqueCC := make([]model.Recipient, 0, len(recipients))
		for i := len(recipients) - 1; i >= 0; i-- {
			cc := recipients[i]
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
		var existing []model.Recipient
		if err := tx.Model(&model.Recipient{}).
			Where("document_id = ?", document.ID).
			Find(&existing).Error; err != nil {
			return err
		}
		existingSet := make(map[uuid.UUID]bool, len(existing))
		for _, cc := range existing {
			existingSet[cc.UserID] = true
		}

		// Prepare sets of IDs to add and remove.
		var toAddIDs []uuid.UUID
		newSet := make(map[uuid.UUID]bool, len(uniqueCC))
		for _, cc := range uniqueCC {
			newSet[cc.UserID] = true
			if !existingSet[cc.UserID] {
				toAddIDs = append(toAddIDs, cc.UserID)
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
			newEntries := make([]model.Recipient, 0, len(toAddIDs))
			for _, uid := range toAddIDs {
				newEntries = append(newEntries, model.Recipient{
					Document: &document,
					UserID:   uid,
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
				Delete(&model.Recipient{}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (t *RecipientRepositoryImpl) Delete(ids []string, documentId string) *helper.ErrorModel {
	var uuids []uuid.UUID

	for _, id := range ids {
		recipientId, errParse := uuid.FromString(id)
		if errParse != nil {
			t.Db.Rollback()
			msg := "Failed to Parse UUID: " + id
			return helper.ErrorCatcher(errParse, 500, &msg)
		}
		uuids = append(uuids, recipientId)
	}

	// Perform a bulk delete using the "IN" condition on the UUID slice
	result := t.Db.Unscoped().Where("user_id IN (?)", uuids).Where("document_id = ?", documentId).Delete(&model.Recipient{})
	if result.Error != nil {
		t.Db.Rollback()
		msg := "Failed to Delete Recipients Data in bulk"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	t.Db.Commit()

	return nil
}
