package controller

import (
	request "Microservice/data/request/Attachment"
	documentAttachment "Microservice/data/response/DocumentAttachment"
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"
	service "Microservice/service/DocumentAttachment"
	userLogService "Microservice/service/UserLog"
	"Microservice/utils"

	"github.com/gin-gonic/gin"
)

type DocumentAttachmentController struct {
	documentAttachmentService service.DocumentAttachmentService
	userLogService            userLogService.UserLogService
}

func NewDocumentAttachmentController(service service.DocumentAttachmentService, userLogService userLogService.UserLogService) *DocumentAttachmentController {
	return &DocumentAttachmentController{documentAttachmentService: service, userLogService: userLogService}
}

func (controller *DocumentAttachmentController) Get(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	stringId := ctx.Param("id")

	documentAttachmentResponse, errDocumentAttachmentResponse := controller.documentAttachmentService.Get(stringId, orgID)
	if errDocumentAttachmentResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentAttachmentResponse)
		return
	}
	utils.SuccessResponse(ctx, documentAttachmentResponse)
}

func (controller *DocumentAttachmentController) GetAll(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	cacheData := utils.GetCache(ctx, "All Attachment", &[]documentAttachment.DocumentAttachmentResponse{})
	if cacheData != nil {
		utils.SuccessResponse(ctx, cacheData)
		return
	}

	documentAttachmentResponse, errDocumentAttachmentResponse := controller.documentAttachmentService.GetAll(orgID)
	if errDocumentAttachmentResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentAttachmentResponse)
		return
	}
	utils.SuccessResponse(ctx, documentAttachmentResponse)
	utils.SetCache(ctx, "All Attachment", documentAttachmentResponse)
}

func (controller *DocumentAttachmentController) Delete(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.AttachmentRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	errDocumentAttachmentResponse := controller.documentAttachmentService.Delete(payload.Id, orgID)
	if errDocumentAttachmentResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentAttachmentResponse)
		return
	}

	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Delete),
			Module: string(enums.DocumentAttachment),
			Log:    helper.ToJSON(payload),
		},
		orgID,
	)

	utils.SuccessResponse(ctx, nil)
}
