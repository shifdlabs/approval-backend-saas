package bookmark

import (
	"Microservice/helper"
	"Microservice/model"

	uuid "github.com/satori/go.uuid"
)

type BookmarkRepository interface {
	AddBookmark(userID, documentID uuid.UUID, orgID string) *helper.ErrorModel
	RemoveBookmark(userID, documentID uuid.UUID, orgID string) *helper.ErrorModel
	IsBookmarked(userID, documentID uuid.UUID, orgID string) (bool, *helper.ErrorModel)
	GetAllBookmarksWithDocuments(userID uuid.UUID, orgID string) ([]model.Document, *helper.ErrorModel)
}
