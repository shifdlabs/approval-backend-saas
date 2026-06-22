package controller

import (
	request "Microservice/data/request/Delegator"
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"
	service "Microservice/service/Delegator"
	userLogService "Microservice/service/UserLog"
	"Microservice/utils"

	"github.com/gin-gonic/gin"
)

type DelegatorController struct {
	delegatorService service.DelegatorService
	userLogService   userLogService.UserLogService
}

func NewDelegatorController(svc service.DelegatorService, userLogSvc userLogService.UserLogService) *DelegatorController {
	return &DelegatorController{delegatorService: svc, userLogService: userLogSvc}
}

func (c *DelegatorController) GetAll(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	ownerID, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	result, err := c.delegatorService.GetAll(*ownerID, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, result)
}

func (c *DelegatorController) Create(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	ownerID, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	var payload request.CreateDelegatorRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if err := c.delegatorService.Create(*ownerID, payload, orgID); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	c.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Create),
		Module: string(enums.Delegator),
	}, orgID)
	utils.SuccessResponse(ctx, nil)
}

func (c *DelegatorController) Update(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	id := ctx.Param("id")

	ownerID, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	var payload request.UpdateDelegatorRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if err := c.delegatorService.Update(id, *ownerID, payload, orgID); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	c.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Update),
		Module: string(enums.Delegator),
	}, orgID)
	utils.SuccessResponse(ctx, nil)
}

func (c *DelegatorController) Delete(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	id := ctx.Param("id")

	ownerID, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	if err := c.delegatorService.Delete(id, *ownerID, orgID); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	c.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Delete),
		Module: string(enums.Delegator),
	}, orgID)
	utils.SuccessResponse(ctx, nil)
}
