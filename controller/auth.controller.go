package controller

import (
	authentication "Microservice/data/request/Authentication"
	response "Microservice/data/response"
	userResponse "Microservice/data/response/User"
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"
	service "Microservice/service/Authentication"
	userService "Microservice/service/User"
	userLogService "Microservice/service/UserLog"
	"Microservice/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService    service.AuthService
	userService    userService.UserService
	userLogService userLogService.UserLogService
}

func NewAuthController(service service.AuthService, userService userService.UserService, userLogSvc userLogService.UserLogService) *AuthController {
	return &AuthController{authService: service, userService: userService, userLogService: userLogSvc}
}

func (controller *AuthController) LogIn(ctx *gin.Context) {
	var payload *authentication.LogInRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure"})
		return
	}

	if errs := helper.ValidateStruct(*payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	loginResult, err := controller.authService.Login(*payload)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	controller.userLogService.CreateLog(model.UserLog{
		UserID: *loginResult.User.ID,
		Action: string(enums.Login),
		Module: string(enums.Authentication),
	}, loginResult.User.OrganizationID.String())

	utils.SuccessResponse(ctx,
		userResponse.LoginResponse{
			AccessToken:      loginResult.AccessToken,
			RefreshToken:     loginResult.RefreshToken,
			UserAbilityRules: controller.GetUserAbilityRules(loginResult.User.Role),
			Id:               loginResult.User.ID.String(),
			Access:           loginResult.User.Access,
			Name:             loginResult.User.FirstName + " " + loginResult.User.LastName,
			Role:             loginResult.User.Role,
			JobPosition:      getPositionName(loginResult.User.Position),
		})
}

func (controller *AuthController) GetUserAbilityRules(userType int) []userResponse.Ability {
	if userType == 99 {
		return []userResponse.Ability{{Action: "manage", Subject: "all"}}
	} else {
		return []userResponse.Ability{{Action: "read", Subject: "all"}}
	}
}

func (controller *AuthController) ForgotPassword(ctx *gin.Context) {
	var payload authentication.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if err := controller.authService.ForgotPassword(payload.Email); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	ctx.JSON(200, response.Response{
		Success: true,
		Code:    200,
		Message: "Password reset email has been sent.",
		Data:    nil,
	})
}

func (controller *AuthController) ResetPassword(ctx *gin.Context) {
	var payload authentication.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure"})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	if err := controller.authService.ResetPassword(payload.Token, payload.NewPassword); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	ctx.JSON(200, response.Response{
		Success: true,
		Code:    200,
		Message: "Password has been successfully reset.",
		Data:    nil,
	})
}

func (controller *AuthController) Logout(ctx *gin.Context) {
	userUUID := helper.GetUserUUID(ctx)
	if userUUID != nil {
		controller.userLogService.CreateLog(model.UserLog{
			UserID: *userUUID,
			Action: string(enums.Logout),
			Module: string(enums.Authentication),
		}, helper.GetOrgID(ctx))
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Code:    200,
		Message: "Success",
		Data:    nil,
	})
}

func getPositionName(position *model.Position) string {
	if position == nil {
		return ""
	}
	return position.Name
}
