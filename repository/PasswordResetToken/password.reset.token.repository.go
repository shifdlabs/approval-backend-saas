package passwordresettoken

import (
	"Microservice/helper"
	"Microservice/model"
)

type PasswordResetTokenRepository interface {
	Create(token model.PasswordResetToken) *helper.ErrorModel
	GetByTokenHash(hash string) (*model.PasswordResetToken, *helper.ErrorModel)
	InvalidateByUserID(userID string) *helper.ErrorModel
	MarkUsed(tokenHash string) *helper.ErrorModel
}
