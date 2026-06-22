package userlog

import (
	response "Microservice/data/response/UserLog"
	repository "Microservice/repository/UserLog"
)

func (t UserLogServiceImpl) mapUserLogToUserLogResponse(rows []repository.UserLogWithName) []response.UserLogResponse {
	responseDocuments := make([]response.UserLogResponse, len(rows))
	for i, row := range rows {
		responseDocuments[i] = response.UserLogResponse{
			Id:       row.ID,
			UserID:   &row.UserID,
			UserName: row.UserName,
			Action:   row.Action,
			Module:   row.Module,
			Log:      row.Log,
			LogDate:  *row.LogDate,
		}
	}
	return responseDocuments
}
