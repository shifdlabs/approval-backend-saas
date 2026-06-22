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
