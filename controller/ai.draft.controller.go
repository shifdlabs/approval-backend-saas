package controller

import (
	request "Microservice/data/request/AIDraft"
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"
	service "Microservice/service/AIDraft"
	userLogService "Microservice/service/UserLog"
	"Microservice/utils"

	"github.com/gin-gonic/gin"
)

type AIDraftController struct {
	service        service.AIDraftService
	userLogService userLogService.UserLogService
}

func NewAIDraftController(svc service.AIDraftService, logSvc userLogService.UserLogService) *AIDraftController {
	return &AIDraftController{service: svc, userLogService: logSvc}
}

// Generate handles POST /api/ai/draft — it asks SAIS for a letter draft and
// returns the text for the editor. Nothing is persisted: the draft only becomes
// a document once the user saves it through the normal document flow.
func (c *AIDraftController) Generate(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.GenerateDraftRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}
	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	// SAIS re-validates the end user's SIS token and logs usage against that
	// user, so the caller's Authorization header is forwarded unchanged.
	result, err := c.service.Generate(ctx.Request.Context(), ctx.GetHeader("Authorization"), payload)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	c.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Create),
		Module: "AI Draft",
	}, orgID)
	utils.SuccessResponse(ctx, result)
}
