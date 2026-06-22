package carbonCopy

import (
	"Microservice/helper"
	"Microservice/model"
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CarbonCopyRepositoryImpl struct {
	Db *gorm.DB
}

func NewCarbonCopyRepositoryImpl(Db *gorm.DB) CarbonCopyRepository {
	return &CarbonCopyRepositoryImpl{Db: Db}
}

func (t *CarbonCopyRepositoryImpl) Create(db gorm.DB, carbonCopys []model.CarbonCopy) *helper.ErrorModel {
	for _, carbonCopy := range carbonCopys {
		if err := db.Create(&carbonCopy).Error; err != nil {
			db.Rollback()
			msg := "Create CarbonCopy Failed"
			return helper.ErrorCatcher(err, 500, &msg)
		}
	}
	return nil
}

func (t *CarbonCopyRepositoryImpl) GetCarbonCopysByDocId(id string) ([]model.CarbonCopy, *helper.ErrorModel) {
	var carbonCopys []model.CarbonCopy
	response := t.Db.Where("document_id = ?", id).Find(&carbonCopys)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "CarbonCopys not found"
		return nil, helper.ErrorCatcher(response.Error, 404, &msg)
	}

	return carbonCopys, nil
}

func (t *CarbonCopyRepositoryImpl) GetAll() ([]model.CarbonCopy, *helper.ErrorModel) {
	var carbonCopys []model.CarbonCopy
	result := t.Db.Find(&carbonCopys)
	if result.Error != nil {
		msg := "CarbonCopys not found"
		return nil, helper.ErrorCatcher(result.Error, 404, &msg)
	}

	return carbonCopys, nil
}

func (t *CarbonCopyRepositoryImpl) Update(document model.Document, carbonCopies []model.CarbonCopy) error {
	// Use a transaction to ensure atomicity:contentReference[oaicite:11]{index=11}.
	return t.Db.Transaction(func(tx *gorm.DB) error {
		// Deduplicate input by UserID (keeping only the last occurrence).
		seen := make(map[uuid.UUID]bool, len(carbonCopies))
		uniqueCC := make([]model.CarbonCopy, 0, len(carbonCopies))
		for i := len(carbonCopies) - 1; i >= 0; i-- {
			cc := carbonCopies[i]
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
		var existing []model.CarbonCopy
		if err := tx.Model(&model.CarbonCopy{}).
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
			newEntries := make([]model.CarbonCopy, 0, len(toAddIDs))
			for _, uid := range toAddIDs {
				newEntries = append(newEntries, model.CarbonCopy{
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
				Delete(&model.CarbonCopy{}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (t *CarbonCopyRepositoryImpl) DeleteCarbonCopy(documentID uint, userID uint) error {
	if userID == 0 {
		return errors.New("invalid userID")
	}
	// Simple delete; no transaction needed for a single operation here.
	return t.Db.Unscoped().Where("document_id = ? AND user_id = ?", documentID, userID).
		Delete(&model.CarbonCopy{}).Error
}
