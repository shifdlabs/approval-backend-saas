package authentication

import (
	model "Microservice/data/model/Authentication"
	authentication "Microservice/data/request/Authentication"
	"Microservice/helper"
)

type AuthService interface {
	Login(payload authentication.LogInRequest) (model.LoginResult, *helper.ErrorModel)
	ForgotPassword(email string) *helper.ErrorModel
	ResetPassword(token, newPassword string) *helper.ErrorModel
}
