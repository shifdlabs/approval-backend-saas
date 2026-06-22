package documenthistory

import (
	"Microservice/helper"
	"Microservice/model"
)

type DocumentHistoryRepository interface {
	Create(document model.DocumentHistory) *helper.ErrorModel
	Get(id string, orgID string) (*model.DocumentHistory, *helper.ErrorModel)
	GetAll(orgID string) ([]model.DocumentHistory, *helper.ErrorModel)
	GetAllHistoryByDocumentId(id string) ([]model.DocumentHistory, *helper.ErrorModel)
	GetLastRejection(id string) (*model.DocumentHistory, *helper.ErrorModel)
	GetLastApprover(id string) (*model.DocumentHistory, *helper.ErrorModel)
	GetHistoriesByAuthorID(userID string, orgID string) ([]model.DocumentHistory, *helper.ErrorModel) // Ganti nama fungsi
	Delete(id string) *helper.ErrorModel
}
