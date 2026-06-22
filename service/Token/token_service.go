package token

import (
	token "Microservice/data/model/Token"
	"Microservice/helper"
)

type TokenService interface {
	RefreshAccessToken(userId string) (token.RefreshTokenResult, *helper.ErrorModel)
}
