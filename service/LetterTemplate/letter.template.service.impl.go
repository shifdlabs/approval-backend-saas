package lettertemplate

import (
	request "Microservice/data/request/LetterTemplate"
	response "Microservice/data/response/LetterTemplate"
	"Microservice/helper"
	"Microservice/model"
	repository "Microservice/repository/LetterTemplate"

	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
)

type LetterTemplateServiceImpl struct {
	Repo     repository.LetterTemplateRepository
	Validate *validator.Validate
}

func NewLetterTemplateServiceImpl(repo repository.LetterTemplateRepository, validate *validator.Validate) LetterTemplateService {
	return &LetterTemplateServiceImpl{Repo: repo, Validate: validate}
}

func toResponse(t model.LetterTemplate) response.LetterTemplateResponse {
	createdAt := ""
	updatedAt := ""
	if t.CreatedAt != nil {
		createdAt = t.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if t.UpdatedAt != nil {
		updatedAt = t.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	return response.LetterTemplateResponse{
		ID:          t.ID.String(),
		Name:        t.Name,
		Description: t.Description,
		Body:        t.Body,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func (s *LetterTemplateServiceImpl) GetAll(orgID string) ([]response.LetterTemplateResponse, *helper.ErrorModel) {
	templates, err := s.Repo.GetAll(orgID)
	if err != nil {
		return nil, err
	}
	var result []response.LetterTemplateResponse
	for _, t := range templates {
		result = append(result, toResponse(t))
	}
	if result == nil {
		result = []response.LetterTemplateResponse{}
	}
	return result, nil
}

func (s *LetterTemplateServiceImpl) GetByID(id string, orgID string) (*response.LetterTemplateResponse, *helper.ErrorModel) {
	t, err := s.Repo.GetByID(id, orgID)
	if err != nil {
		return nil, err
	}
	r := toResponse(*t)
	return &r, nil
}

func (s *LetterTemplateServiceImpl) Create(req request.CreateLetterTemplateRequest, orgID string) (*response.LetterTemplateResponse, *helper.ErrorModel) {
	if errV := s.Validate.Struct(req); errV != nil {
		msg := "Validation error"
		return nil, helper.ErrorCatcher(errV, 400, &msg)
	}

	orgUUID, errParse := uuid.FromString(orgID)
	if errParse != nil {
		msg := "Invalid Organization ID"
		return nil, helper.ErrorCatcher(errParse, 500, &msg)
	}

	t, err := s.Repo.Create(model.LetterTemplate{
		OrganizationID: &orgUUID,
		Name:           req.Name,
		Description:    req.Description,
		Body:           req.Body,
	})
	if err != nil {
		return nil, err
	}
	r := toResponse(*t)
	return &r, nil
}

func (s *LetterTemplateServiceImpl) Update(id string, req request.UpdateLetterTemplateRequest, orgID string) (*response.LetterTemplateResponse, *helper.ErrorModel) {
	if errV := s.Validate.Struct(req); errV != nil {
		msg := "Validation error"
		return nil, helper.ErrorCatcher(errV, 400, &msg)
	}
	t, err := s.Repo.Update(id, model.LetterTemplate{
		Name:        req.Name,
		Description: req.Description,
		Body:        req.Body,
	}, orgID)
	if err != nil {
		return nil, err
	}
	r := toResponse(*t)
	return &r, nil
}

func (s *LetterTemplateServiceImpl) Delete(id string, orgID string) *helper.ErrorModel {
	return s.Repo.Delete(id, orgID)
}
