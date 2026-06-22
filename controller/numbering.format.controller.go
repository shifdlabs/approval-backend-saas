package controller

import (
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"
	"Microservice/utils"

	request "Microservice/data/request/NumberingFormat"
	service "Microservice/service/NumberingFormat"
	userLogService "Microservice/service/UserLog"

	"github.com/gin-gonic/gin"
)

type NumberingFormatController struct {
	numberingFormatService service.NumberingFormatService
	userLogService         userLogService.UserLogService
}

func NewNumberingFormatController(service service.NumberingFormatService, userLogService userLogService.UserLogService) *NumberingFormatController {
	return &NumberingFormatController{numberingFormatService: service, userLogService: userLogService}
}

func (controller *NumberingFormatController) Create(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.NumberingFormatRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	err := controller.numberingFormatService.Create(payload, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Create),
			Module: string(enums.NumberingFormat),
			Log:    helper.ToJSON(payload),
		},
		orgID,
	)

	utils.SuccessResponse(ctx, nil)
}

func (controller *NumberingFormatController) GetAll(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	documentSequenceResponse, errDocumentSequenceResponse := controller.numberingFormatService.GetAll(orgID)
	if errDocumentSequenceResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentSequenceResponse)
		return
	}
	utils.SuccessResponse(ctx, documentSequenceResponse)
}

func (controller *NumberingFormatController) GetAllWithGrouped(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	documentSequenceResponse, errDocumentSequenceResponse := controller.numberingFormatService.GetAllWithGrouped(orgID)
	if errDocumentSequenceResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentSequenceResponse)
		return
	}
	utils.SuccessResponse(ctx, documentSequenceResponse)
}

func (controller *NumberingFormatController) Delete(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	stringID := ctx.Param("id")
	errResponse := controller.numberingFormatService.Delete(stringID, orgID)
	if errResponse != nil {
		utils.ErrorResponse(ctx, *errResponse)
		return
	}

	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Delete),
			Module: string(enums.NumberingGroup),
		},
		orgID,
	)

	utils.SuccessResponse(ctx, nil)
}
