package documentnumbers

import (
	request "Microservice/data/request/DocumentNumbers"
	response "Microservice/data/response/DocumentNumbers"
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"

	uuid "github.com/satori/go.uuid"
)

type DocumentNumbersService interface {
	Create(request request.DocumentNumbersRequest, userId string, document *model.Document, state enums.DocumentNumberState, orgID string) *helper.ErrorModel
	Update(id string, document *model.Document, state enums.DocumentNumberState, orgID string) *helper.ErrorModel
	GetAll(orgID string) ([]response.DocumentNumbersResponse, *helper.ErrorModel)
	Get(id string, orgID string) (*response.DocumentNumbersResponse, *helper.ErrorModel)
	GetByDocumentID(id uuid.UUID, orgID string) (*response.DocumentNumbersResponse, *helper.ErrorModel)
	GetAllByUserId(userId string, orgID string) ([]response.DocumentNumbersResponse, *helper.ErrorModel)
	Delete(id string, orgID string) *helper.ErrorModel
}
