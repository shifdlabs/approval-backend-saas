package bookmark

import (
	request "Microservice/data/request/Bookmark"
	"Microservice/helper"
	"Microservice/model"
	bookmarkRepository "Microservice/repository/Bookmark"

	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
)

type BookmarkServiceImpl struct {
	BookmarkRepository bookmarkRepository.BookmarkRepository
	Validate           *validator.Validate
}

func NewBookmarkServiceImpl(
	bookmarkRepository bookmarkRepository.BookmarkRepository,
	validate *validator.Validate) BookmarkService {
	return &BookmarkServiceImpl{
		BookmarkRepository: bookmarkRepository,
		Validate:           validate,
	}
}

// AddBookmark menambahkan bookmark baru
func (s BookmarkServiceImpl) AddBookmark(request request.BookmarkRequest, orgID string) *helper.ErrorModel {
	// Validasi request
	errStructure := s.Validate.Struct(request)
	if errStructure != nil {
		msg := "Structure Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	// Parse UUID
	userID, errUser := uuid.FromString(request.UserID)
	if errUser != nil {
		msg := "Invalid UserID"
		return helper.ErrorCatcher(errUser, 400, &msg)
	}

	documentID, errDoc := uuid.FromString(request.DocumentID)
	if errDoc != nil {
		msg := "Invalid DocumentID"
		return helper.ErrorCatcher(errDoc, 400, &msg)
	}

	// Tambahkan bookmark
	errCreate := s.BookmarkRepository.AddBookmark(userID, documentID, orgID)
	if errCreate != nil {
		msg := "Failed to add bookmark"
		return helper.ErrorCatcher(errCreate, 500, &msg)
	}

	return nil
}

// RemoveBookmark menghapus bookmark
func (s BookmarkServiceImpl) RemoveBookmark(request request.BookmarkRequest, orgID string) *helper.ErrorModel {
	// Validasi request
	errStructure := s.Validate.Struct(request)
	if errStructure != nil {
		msg := "Structure Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	// Parse UUID
	userID, errUser := uuid.FromString(request.UserID)
	if errUser != nil {
		msg := "Invalid UserID"
		return helper.ErrorCatcher(errUser, 400, &msg)
	}

	documentID, errDoc := uuid.FromString(request.DocumentID)
	if errDoc != nil {
		msg := "Invalid DocumentID"
		return helper.ErrorCatcher(errDoc, 400, &msg)
	}

	// Hapus bookmark
	errRemove := s.BookmarkRepository.RemoveBookmark(userID, documentID, orgID)
	if errRemove != nil {
		msg := "Failed to remove bookmark"
		return helper.ErrorCatcher(errRemove, 500, &msg)
	}

	return nil
}

// IsBookmarked memeriksa apakah dokumen sudah di-bookmark
func (s BookmarkServiceImpl) IsBookmarked(request request.BookmarkRequest, orgID string) (bool, *helper.ErrorModel) {
	// Validasi request
	errStructure := s.Validate.Struct(request)
	if errStructure != nil {
		msg := "Structure Error"
		return false, helper.ErrorCatcher(errStructure, 500, &msg)
	}

	// Parse UUID
	userID, errUser := uuid.FromString(request.UserID)
	if errUser != nil {
		msg := "Invalid UserID"
		return false, helper.ErrorCatcher(errUser, 400, &msg)
	}

	documentID, errDoc := uuid.FromString(request.DocumentID)
	if errDoc != nil {
		msg := "Invalid DocumentID"
		return false, helper.ErrorCatcher(errDoc, 400, &msg)
	}

	// Periksa bookmark
	isBookmarked, errCheck := s.BookmarkRepository.IsBookmarked(userID, documentID, orgID)
	if errCheck != nil {
		msg := "Failed to check bookmark status"
		return false, helper.ErrorCatcher(errCheck, 500, &msg)
	}

	return isBookmarked, nil
}

func (s BookmarkServiceImpl) GetAllBookmarksWithDocuments(userID string, orgID string) ([]model.Document, *helper.ErrorModel) {
	// Parse UUID
	parsedUserID, err := uuid.FromString(userID)
	if err != nil {
		msg := "Invalid UserID"
		return nil, helper.ErrorCatcher(err, 400, &msg)
	}

	// Panggil repository
	documents, errFetch := s.BookmarkRepository.GetAllBookmarksWithDocuments(parsedUserID, orgID)
	if errFetch != nil {
		return nil, errFetch
	}

	return documents, nil
}
