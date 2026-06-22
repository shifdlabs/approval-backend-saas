package controller

import (
	"Microservice/helper"
	service "Microservice/service/UserLog"
	"Microservice/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

type UserLogController struct {
	userLogService service.UserLogService
}

func NewUserLogController(service service.UserLogService) *UserLogController {
	return &UserLogController{userLogService: service}
}

func (controller *UserLogController) GetAll(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	userLogResponse, errUserLogResponse := controller.userLogService.GetAll(orgID)
	if errUserLogResponse != nil {
		utils.ErrorResponse(ctx, *errUserLogResponse)
		return
	}
	utils.SuccessResponse(ctx, userLogResponse)
}

func (controller *UserLogController) Export(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	data, err := controller.userLogService.Export(orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	ctx.Header("Content-Disposition", `attachment; filename="activity-log.xlsx"`)
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(data)))
	ctx.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}
