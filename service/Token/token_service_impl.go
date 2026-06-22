package token

import (
	"Microservice/config"
	token "Microservice/data/model/Token"
	"Microservice/helper"
	repository "Microservice/repository/User"
)

type TokenServiceImpl struct {
	UserRepository repository.UserRepository
}

func NewTokenServiceImpl(userRepository repository.UserRepository) TokenService {
	return &TokenServiceImpl{
		UserRepository: userRepository,
	}
}

func (t TokenServiceImpl) RefreshAccessToken(userId string) (token.RefreshTokenResult, *helper.ErrorModel) {
	// Unrouted in Phase 2 (SIS owns refresh) — unscoped lookup since no org_id is available here.
	user, err := t.UserRepository.GetUnscoped(userId, true)
	if err != nil {
		return token.RefreshTokenResult{
			AccessToken:  nil,
			RefreshToken: nil,
		}, err
	}

	env, _ := config.LoadConfig(".")

	accessTokenDetails, accessTokenErr := config.CreateAccessToken(user, env.AccessTokenExpiresIn, env.AccessTokenPrivateKey)
	if accessTokenErr != nil {
		msg := "Failed create access token"
		return token.RefreshTokenResult{
			AccessToken:  nil,
			RefreshToken: nil,
		}, helper.ErrorCatcher(accessTokenErr, 500, &msg)
	}

	refreshTokenDetails, refreshTokenErr := config.CreateRefreshToken(user, env.RefreshTokenExpiresIn, env.RefreshTokenPrivateKey)
	if refreshTokenErr != nil {
		msg := "Failed create refresh token"
		return token.RefreshTokenResult{
			AccessToken:  nil,
			RefreshToken: nil,
		}, helper.ErrorCatcher(accessTokenErr, 500, &msg)
	}

	return token.RefreshTokenResult{
		AccessToken:  accessTokenDetails,
		RefreshToken: refreshTokenDetails,
	}, nil
}
