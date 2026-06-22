package controller

import (
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"
	"Microservice/utils"

	request "Microservice/data/request/Position"
	service "Microservice/service/Position"
	userLogService "Microservice/service/UserLog"

	"github.com/gin-gonic/gin"
)

type PositionController struct {
	positionService service.PositionService
	userLogService  userLogService.UserLogService
}

func NewPositionController(service service.PositionService, userLogService userLogService.UserLogService) *PositionController {
	return &PositionController{positionService: service, userLogService: userLogService}
}

func (controller *PositionController) Get(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	stringID := ctx.Param("id")

	positionResponse, err := controller.positionService.Get(stringID, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, positionResponse)
}

func (controller *PositionController) GetAll(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	positionResponse, err := controller.positionService.GetAll(orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, positionResponse)
}

func (controller *PositionController) Create(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.CreatePositionRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	err := controller.positionService.Create(payload, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Create),
			Module: string(enums.Position),
			Log:    helper.ToJSON(payload),
		},
		orgID,
	)

	utils.SuccessResponse(ctx, nil)
}

func (controller *PositionController) Update(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.UpdatePositionRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	err := controller.positionService.Update(payload, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Update),
			Module: string(enums.Position),
			Log:    helper.ToJSON(payload),
		},
		orgID,
	)

	utils.SuccessResponse(ctx, nil)
}

func (controller *PositionController) Delete(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	stringID := ctx.Param("id")
	errResponse := controller.positionService.Delete(stringID, orgID)
	if errResponse != nil {
		utils.ErrorResponse(ctx, *errResponse)
		return
	}

	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Delete),
			Module: string(enums.Position),
		},
		orgID,
	)

	utils.SuccessResponse(ctx, nil)
}
