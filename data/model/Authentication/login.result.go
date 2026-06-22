package authentication

import (
	"Microservice/model"
)

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	User         *model.User
}
