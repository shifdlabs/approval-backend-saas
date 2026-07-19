package controller

import (
	"Microservice/helper"
	"Microservice/utils"

	request "Microservice/data/request/Bookmark"
	service "Microservice/service/Bookmark"

	"github.com/gin-gonic/gin"
)

type BookmarkController struct {
	bookmarkService service.BookmarkService
}

func NewBookmarkController(service service.BookmarkService) *BookmarkController {
	return &BookmarkController{bookmarkService: service}
}

// bindCallerBookmark overwrites the request UserID with the authenticated
// user's id (JWT sub), so a member cannot bookmark on another member's behalf
// by forging the userId in the body (AUDIT IDOR). Returns false (and writes the
// error response) if the caller id is unavailable.
func bindCallerBookmark(ctx *gin.Context, payload *request.BookmarkRequest) bool {
	callerID, errCaller := helper.GetUserId(ctx)
	if errCaller != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return false
	}
	payload.UserID = *callerID
	return true
}

func (controller *BookmarkController) AddBookmarkHandler(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.BookmarkRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	// Bind the bookmark to the authenticated user, not the client-supplied id,
	// so a member cannot manage another member's bookmarks (AUDIT IDOR).
	if !bindCallerBookmark(ctx, &payload) {
		return
	}

	err := controller.bookmarkService.AddBookmark(payload, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, nil)
}

func (controller *BookmarkController) RemoveBookmarkHandler(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.BookmarkRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if !bindCallerBookmark(ctx, &payload) {
		return
	}

	err := controller.bookmarkService.RemoveBookmark(payload, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, nil)
}

func (controller *BookmarkController) IsBookmarkedHandler(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.BookmarkRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if !bindCallerBookmark(ctx, &payload) {
		return
	}

	isBookmarked, err := controller.bookmarkService.IsBookmarked(payload, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, gin.H{"isBookmarked": isBookmarked})
}

func (controller *BookmarkController) GetAllBookmarksWithDocumentsHandler(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	documents, err := controller.bookmarkService.GetAllBookmarksWithDocuments(*id, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, documents)
}
