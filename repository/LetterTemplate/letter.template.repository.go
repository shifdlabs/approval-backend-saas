package lettertemplate

import (
	"Microservice/helper"
	"Microservice/model"
)

type LetterTemplateRepository interface {
	GetAll(orgID string) ([]model.LetterTemplate, *helper.ErrorModel)
	GetByID(id string, orgID string) (*model.LetterTemplate, *helper.ErrorModel)
	Create(template model.LetterTemplate) (*model.LetterTemplate, *helper.ErrorModel)
	Update(id string, template model.LetterTemplate, orgID string) (*model.LetterTemplate, *helper.ErrorModel)
	Delete(id string, orgID string) *helper.ErrorModel
}
