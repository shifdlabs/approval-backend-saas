package bookmark

import (
	request "Microservice/data/request/Bookmark"
	"Microservice/helper"
	"Microservice/model"
)

type BookmarkService interface {
	AddBookmark(request request.BookmarkRequest, orgID string) *helper.ErrorModel
	RemoveBookmark(request request.BookmarkRequest, orgID string) *helper.ErrorModel
	IsBookmarked(request request.BookmarkRequest, orgID string) (bool, *helper.ErrorModel)
	GetAllBookmarksWithDocuments(userID string, orgID string) ([]model.Document, *helper.ErrorModel)
}
