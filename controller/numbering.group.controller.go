package controller

import (
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"
	"Microservice/utils"

	request "Microservice/data/request/NumberingGroup"
	service "Microservice/service/NumberingGroup"
	userLogService "Microservice/service/UserLog"

	"github.com/gin-gonic/gin"
)

type NumberingGroupController struct {
	numberingGroupService service.NumberingGroupService
	userLogService        userLogService.UserLogService
}

func NewNumberingGroupController(service service.NumberingGroupService, userLogService userLogService.UserLogService) *NumberingGroupController {
	return &NumberingGroupController{numberingGroupService: service, userLogService: userLogService}
}

func (controller *NumberingGroupController) Get(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	stringID := ctx.Param("id")

	numberingGroupResponse, err := controller.numberingGroupService.Get(stringID, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, numberingGroupResponse)
}

func (controller *NumberingGroupController) GetAll(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	numberingGroupResponse, err := controller.numberingGroupService.GetAll(orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, numberingGroupResponse)
}

func (controller *NumberingGroupController) Create(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.NumberingGroupRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	err := controller.numberingGroupService.Create(payload, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Create),
			Module: string(enums.NumberingGroup),
			Log:    helper.ToJSON(payload),
		},
		orgID,
	)

	utils.SuccessResponse(ctx, nil)
}

func (controller *NumberingGroupController) Delete(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	stringID := ctx.Param("id")
	errResponse := controller.numberingGroupService.Delete(stringID, orgID)
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
