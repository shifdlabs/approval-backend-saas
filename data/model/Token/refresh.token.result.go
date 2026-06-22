package token

import (
	"Microservice/config"
)

type RefreshTokenResult struct {
	AccessToken  *config.TokenDetails
	RefreshToken *config.TokenDetails
}
