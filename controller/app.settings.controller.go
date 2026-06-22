package controller

import (
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"
	"Microservice/utils"

	request "Microservice/data/request/AppSettings"
	service "Microservice/service/AppSettings"
	userLogService "Microservice/service/UserLog"

	"github.com/gin-gonic/gin"
)

type AppSettingsController struct {
	appSettingsService service.AppSettingService
	userLogService     userLogService.UserLogService
}

func NewAppSettingsController(service service.AppSettingService, userLogSvc userLogService.UserLogService) *AppSettingsController {
	return &AppSettingsController{appSettingsService: service, userLogService: userLogSvc}
}

func (controller *AppSettingsController) GetAll(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	appSettingsResponse, err := controller.appSettingsService.GetAll(orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, appSettingsResponse)
}

func (controller *AppSettingsController) Update(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.AppSettingRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	err := controller.appSettingsService.Update(payload, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	controller.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Update),
		Module: string(enums.AppSettingsModule),
	}, orgID)
	utils.SuccessResponse(ctx, nil)
}
