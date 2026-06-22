package controller

import (
	"Microservice/helper"
	service "Microservice/service/DocumentSequence"
	"Microservice/utils"

	"github.com/gin-gonic/gin"
)

type DocumentSequenceController struct {
	documentSequenceService service.DocumentSequenceService
}

func NewDocumentSequenceController(service service.DocumentSequenceService) *DocumentSequenceController {
	return &DocumentSequenceController{documentSequenceService: service}
}

func (controller *DocumentSequenceController) Get(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	stringId := ctx.Param("id")

	documentSequenceResponse, errDocumentSequenceResponse := controller.documentSequenceService.Get(stringId, orgID)
	if errDocumentSequenceResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentSequenceResponse)
		return
	}
	utils.SuccessResponse(ctx, documentSequenceResponse)
}

func (controller *DocumentSequenceController) GetProgress(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	documentSequenceResponses, err := controller.documentSequenceService.GetProgressByAuthorID(*id, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, documentSequenceResponses)
}
