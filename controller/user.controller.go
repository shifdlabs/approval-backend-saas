package controller

import (
	request "Microservice/data/request/User"
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"
	service "Microservice/service/User"
	userLogService "Microservice/service/UserLog"
	"Microservice/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

const maxUploadSize = 10 * 1024 * 1024 // 10 MB

type UserController struct {
	userService    service.UserService
	userLogService userLogService.UserLogService
}

func NewUserController(service service.UserService, userLogSvc userLogService.UserLogService) *UserController {
	return &UserController{userService: service, userLogService: userLogSvc}
}

func (controller *UserController) Get(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	stringID, errorParseToken := helper.GetUserId(ctx)
	if errorParseToken != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	userResponse, errResponse := controller.userService.Get(*stringID, orgID)
	if errResponse != nil {
		utils.ErrorResponse(ctx, *errResponse)
		return
	}
	utils.SuccessResponse(ctx, userResponse)
}

func (controller *UserController) GetUserByID(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	stringID := ctx.Param("id")
	if stringID == "" {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	userResponse, errResponse := controller.userService.Get(stringID, orgID)
	if errResponse != nil {
		utils.ErrorResponse(ctx, *errResponse)
		return
	}
	utils.SuccessResponse(ctx, userResponse)
}

func (controller *UserController) GetAll(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	userResponse, errResponse := controller.userService.GetAll(orgID)
	if errResponse != nil {
		utils.ErrorResponse(ctx, *errResponse)
		return
	}
	utils.SuccessResponse(ctx, userResponse)
}

func (controller *UserController) GetAllUserExceptCurrent(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	stringID, errorParseToken := helper.GetUserId(ctx)
	if errorParseToken != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	userResponse, errResponse := controller.userService.GetAllUserExceptCurrent(*stringID, orgID)
	if errResponse != nil {
		utils.ErrorResponse(ctx, *errResponse)
		return
	}
	utils.SuccessResponse(ctx, userResponse)
}

func (controller *UserController) Create(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.CreateUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if err := controller.userService.Create(payload, orgID); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	controller.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Create),
		Module: string(enums.User),
	}, orgID)
	utils.SuccessResponse(ctx, nil)
}

func (controller *UserController) Update(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if err := controller.userService.Update(payload, orgID); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	controller.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Update),
		Module: string(enums.User),
	}, orgID)
	utils.SuccessResponse(ctx, nil)
}

func (controller *UserController) Delete(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	stringId := ctx.Param("id")

	if err := controller.userService.Delete(stringId, orgID); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	controller.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Delete),
		Module: string(enums.User),
	}, orgID)
	utils.SuccessResponse(ctx, nil)
}

func (controller *UserController) MultipleDelete(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.DeleteMultipleUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if err := controller.userService.MultipleDelete(payload.IDs, orgID); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	controller.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Delete),
		Module: string(enums.User),
	}, orgID)
	utils.SuccessResponse(ctx, nil)
}

func (controller *UserController) UpdateEmail(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	stringID, errorParseToken := helper.GetUserId(ctx)
	if errorParseToken != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	var payload request.UpdateEmailRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if err := controller.userService.UpdateEmail(*stringID, payload, orgID); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	controller.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Update),
		Module: string(enums.Profile),
	}, orgID)
	utils.SuccessResponse(ctx, nil)
}

func (controller *UserController) UpdateBiodata(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	stringID, errorParseToken := helper.GetUserId(ctx)
	if errorParseToken != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	var payload request.UpdateBiodataRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if err := controller.userService.UpdateBiodata(*stringID, payload, orgID); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	controller.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Update),
		Module: string(enums.Profile),
	}, orgID)
	utils.SuccessResponse(ctx, nil)
}

func (controller *UserController) UpdateRole(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if err := controller.userService.UpdateRole(payload, orgID); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	controller.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Update),
		Module: string(enums.User),
	}, orgID)
	utils.SuccessResponse(ctx, nil)
}

func (controller *UserController) UpdatePassword(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.UpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if err := controller.userService.UpdatePassword(payload, orgID); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	controller.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Update),
		Module: string(enums.Profile),
	}, orgID)
	utils.SuccessResponse(ctx, nil)
}

func (controller *UserController) UpdateAccess(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.UpdateAccessRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if err := controller.userService.UpdateAccess(payload, orgID); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	controller.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(enums.Update),
		Module: string(enums.User),
	}, orgID)
	utils.SuccessResponse(ctx, nil)
}

func (controller *UserController) PreviewImport(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "File is required"})
		return
	}

	if file.Size > maxUploadSize {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "File size must not exceed 10 MB"})
		return
	}

	if !strings.HasSuffix(strings.ToLower(file.Filename), ".xlsx") &&
		!strings.HasSuffix(strings.ToLower(file.Filename), ".xls") {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Only Excel files (.xlsx, .xls) are allowed"})
		return
	}

	columnMappingJSON := ctx.PostForm("columnMapping")
	if columnMappingJSON == "" {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Column mapping is required"})
		return
	}

	previewResponse, errResponse := controller.userService.PreviewImport(file, columnMappingJSON, orgID)
	if errResponse != nil {
		utils.ErrorResponse(ctx, *errResponse)
		return
	}
	utils.SuccessResponse(ctx, previewResponse)
}

func (controller *UserController) UnlockUser(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	userId := ctx.Param("userId")
	if userId == "" {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "User ID is required"})
		return
	}

	if err := controller.userService.UnlockUser(userId, orgID); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, nil)
}

func (controller *UserController) BulkImport(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.BulkImportUsersRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	importResponse, errResponse := controller.userService.BulkImport(payload, orgID)
	if errResponse != nil {
		utils.ErrorResponse(ctx, *errResponse)
		return
	}
	utils.SuccessResponse(ctx, importResponse)
}
