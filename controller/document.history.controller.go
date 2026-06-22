package controller

import (
	documentHistory "Microservice/data/response/DocumentHistory"
	"Microservice/helper"
	service "Microservice/service/DocumentHistory"
	"Microservice/utils"

	"github.com/gin-gonic/gin"
)

type DocumentHistoryController struct {
	documentHistoryService service.DocumentHistoryService
}

func NewDocumentHistoryController(service service.DocumentHistoryService) *DocumentHistoryController {
	return &DocumentHistoryController{documentHistoryService: service}
}

func (controller *DocumentHistoryController) Get(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	stringId := ctx.Param("id")

	documentHistoryResponse, errDocumentHistoryResponse := controller.documentHistoryService.Get(stringId, orgID)
	if errDocumentHistoryResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentHistoryResponse)
		return
	}
	utils.SuccessResponse(ctx, documentHistoryResponse)
}

func (controller *DocumentHistoryController) GetAll(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	cacheData := utils.GetCache(ctx, "All History", &[]documentHistory.DocumentHistoryResponse{})
	if cacheData != nil {
		utils.SuccessResponse(ctx, cacheData)
		return
	}

	documentHistoryResponse, errDocumentHistoryResponse := controller.documentHistoryService.GetAll(orgID)
	if errDocumentHistoryResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentHistoryResponse)
		return
	}
	utils.SuccessResponse(ctx, documentHistoryResponse)
	utils.SetCache(ctx, "All History", documentHistoryResponse)
}

func (controller *DocumentHistoryController) GetRejectedWithDocumentAndUser(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	documentHistoryResponse, err := controller.documentHistoryService.FetchHistoriesByUserID(*id, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, documentHistoryResponse)
}
