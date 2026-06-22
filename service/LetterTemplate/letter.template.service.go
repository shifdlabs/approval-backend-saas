package lettertemplate

import (
	request "Microservice/data/request/LetterTemplate"
	response "Microservice/data/response/LetterTemplate"
	"Microservice/helper"
)

type LetterTemplateService interface {
	GetAll(orgID string) ([]response.LetterTemplateResponse, *helper.ErrorModel)
	GetByID(id string, orgID string) (*response.LetterTemplateResponse, *helper.ErrorModel)
	Create(req request.CreateLetterTemplateRequest, orgID string) (*response.LetterTemplateResponse, *helper.ErrorModel)
	Update(id string, req request.UpdateLetterTemplateRequest, orgID string) (*response.LetterTemplateResponse, *helper.ErrorModel)
	Delete(id string, orgID string) *helper.ErrorModel
}
