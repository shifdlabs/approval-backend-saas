package controller

import (
	request "Microservice/data/request/LetterTemplate"
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"
	service "Microservice/service/LetterTemplate"
	userLogService "Microservice/service/UserLog"
	"Microservice/utils"

	"github.com/gin-gonic/gin"
)

type LetterTemplateController struct {
	service        service.LetterTemplateService
	userLogService userLogService.UserLogService
}

func NewLetterTemplateController(svc service.LetterTemplateService, logSvc userLogService.UserLogService) *LetterTemplateController {
	return &LetterTemplateController{service: svc, userLogService: logSvc}
}

func (c *LetterTemplateController) GetAll(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	templates, err := c.service.GetAll(orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, templates)
}

func (c *LetterTemplateController) GetByID(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	id := ctx.Param("id")
	template, err := c.service.GetByID(id, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, template)
}

func (c *LetterTemplateController) Create(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.CreateLetterTemplateRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}
	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}
	result, err := c.service.Create(payload, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	c.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Create),
		Module: "Letter Template",
	}, orgID)
	utils.SuccessResponse(ctx, result)
}

func (c *LetterTemplateController) Update(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	id := ctx.Param("id")
	var payload request.UpdateLetterTemplateRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}
	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}
	result, err := c.service.Update(id, payload, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	c.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Update),
		Module: "Letter Template",
	}, orgID)
	utils.SuccessResponse(ctx, result)
}

func (c *LetterTemplateController) Delete(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	id := ctx.Param("id")
	if err := c.service.Delete(id, orgID); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	c.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Delete),
		Module: "Letter Template",
	}, orgID)
	utils.SuccessResponse(ctx, nil)
}
