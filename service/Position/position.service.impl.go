package position

import (
	request "Microservice/data/request/Position"
	response "Microservice/data/response/Position"
	"Microservice/helper"
	"Microservice/model"
	repository "Microservice/repository/Position"

	uuid "github.com/satori/go.uuid"

	"github.com/go-playground/validator/v10"
)

type PositionServiceImpl struct {
	PositionRepository repository.PositionRepository
	Validate           *validator.Validate
}

func NewPositionServiceImpl(
	reportRepository repository.PositionRepository,
	validate *validator.Validate) PositionService {
	return &PositionServiceImpl{
		PositionRepository: reportRepository,
		Validate:           validate,
	}
}

func (t PositionServiceImpl) Create(position request.CreatePositionRequest, orgID string) *helper.ErrorModel {
	errStructure := t.Validate.Struct(position)

	if errStructure != nil {
		msg := "Structure Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	orgUUID, errParse := uuid.FromString(orgID)
	if errParse != nil {
		msg := "Invalid Organization ID"
		return helper.ErrorCatcher(errParse, 500, &msg)
	}

	errCreate := t.PositionRepository.Create(model.Position{Name: position.Name, OrganizationID: &orgUUID})
	if errCreate != nil {
		return errCreate
	}

	return nil
}

func (t PositionServiceImpl) Get(id string, orgID string) (*response.PositionResponse, *helper.ErrorModel) {
	report, errFetch := t.PositionRepository.Get(id, orgID)

	if errFetch != nil {
		return nil, errFetch
	} else if report == nil {
		return nil, nil
	}

	reportResponse := t.convertPositionToPositionResponse(*report)
	return &reportResponse, nil
}

func (t PositionServiceImpl) GetAll(orgID string) ([]response.PositionResponse, *helper.ErrorModel) {
	result, errFetch := t.PositionRepository.GetAll(orgID)

	if errFetch != nil {
		return nil, errFetch
	} else {
		return t.mapPositionToPositionResponse(result), nil
	}
}

func (t PositionServiceImpl) Update(position request.UpdatePositionRequest, orgID string) *helper.ErrorModel {
	errStructure := t.Validate.Struct(position)
	if errStructure != nil {
		msg := "Structure Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	id, errParse := uuid.FromString(position.ID)
	if errParse != nil {
		msg := "Parse Error"
		return helper.ErrorCatcher(errParse, 500, &msg)
	}

	errUpdate := t.PositionRepository.Update(model.Position{ID: &id, Name: position.Name}, orgID)
	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

func (t PositionServiceImpl) Delete(id string, orgID string) *helper.ErrorModel {
	errResponse := t.PositionRepository.Delete(id, orgID)

	if errResponse != nil {
		return errResponse
	}

	return nil
}
