package userlog

import (
	response "Microservice/data/response/UserLog"
	"Microservice/helper"
	"Microservice/model"
)

type UserLogService interface {
	GetAll(orgID string) ([]response.UserLogResponse, *helper.ErrorModel)
	CreateLog(log model.UserLog, orgID string)
	Export(orgID string) ([]byte, *helper.ErrorModel)
}
