package failedloginattempt

import (
	"Microservice/helper"
	"Microservice/model"
)

type FailedLoginAttemptRepository interface {
	Create(attempt model.FailedLoginAttempt) *helper.ErrorModel
	CountByUserId(userId string) (int64, *helper.ErrorModel)
	DeleteByUserId(userId string) *helper.ErrorModel
}
