package documentHistory

import (
	response "Microservice/data/response/DocumentHistory"
	"Microservice/helper"
)

type DocumentHistoryService interface {
	Get(id string, orgID string) (*response.DocumentHistoryResponse, *helper.ErrorModel)
	GetAll(orgID string) ([]response.DocumentHistoryResponse, *helper.ErrorModel)
	Delete(id string) *helper.ErrorModel
	FetchHistoriesByUserID(userID string, orgID string) ([]response.DocumentHistoryResponse, *helper.ErrorModel)
}
