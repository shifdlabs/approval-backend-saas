package carbonCopy

import (
	request "Microservice/data/request/CarbonCopy"
	"Microservice/helper"
	"Microservice/model"
	repository "Microservice/repository/CarbonCopy"
	documentRepository "Microservice/repository/Document"

	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type CarbonCopyServiceImpl struct {
	CarbonCopyRepository repository.CarbonCopyRepository
	DocumentRepository   documentRepository.DocumentRepository
	Validate             *validator.Validate
	Db                   *gorm.DB
}

func NewCarbonCopyServiceImpl(
	carbonCopyRepository repository.CarbonCopyRepository,
	documentRepsitory documentRepository.DocumentRepository,
	Db *gorm.DB,
	validate *validator.Validate) CarbonCopyService {
	return &CarbonCopyServiceImpl{
		CarbonCopyRepository: carbonCopyRepository,
		DocumentRepository:   documentRepsitory,
		Db:                   Db,
		Validate:             validate,
	}
}

func (t CarbonCopyServiceImpl) Create(request request.CarbonCopyRequest, orgID string) *helper.ErrorModel {
	trx := t.Db.Begin()

	errStructure := t.Validate.Struct(request)
	if errStructure != nil {
		trx.Rollback()
		msg := "Structure Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	document, errFetch := t.DocumentRepository.Get(request.DocumentId, orgID)
	if errFetch != nil {
		trx.Rollback()
		msg := "Structure Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	var carbonCopys []model.CarbonCopy
	for _, userId := range request.UserIds {
		userUuid, errorParse := uuid.FromString(userId)
		if errorParse != nil {
			msg := "Structure Error"
			return helper.ErrorCatcher(errStructure, 500, &msg)
		}

		carbonCopys = append(carbonCopys, model.CarbonCopy{
			Document: document,
			UserID:   userUuid,
		})
	}

	errCreate := t.CarbonCopyRepository.Create(*trx, carbonCopys)
	if errCreate != nil {
		trx.Rollback()
		msg := "Create Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	return nil
}

func (t CarbonCopyServiceImpl) Update(request request.CarbonCopyRequest, orgID string) *helper.ErrorModel {
	errStructure := t.Validate.Struct(request)
	if errStructure != nil {
		msg := "Structure Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	document, errFetch := t.DocumentRepository.Get(request.DocumentId, orgID)
	if errFetch != nil {
		msg := "Structure Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	var carbonCopys []model.CarbonCopy
	for _, userId := range request.UserIds {
		userUuid, errorParse := uuid.FromString(userId)
		if errorParse != nil {
			msg := "Structure Error"
			return helper.ErrorCatcher(errStructure, 500, &msg)
		}

		carbonCopys = append(carbonCopys, model.CarbonCopy{
			Document: document,
			UserID:   userUuid,
		})
	}

	errCreate := t.CarbonCopyRepository.Update(*document, carbonCopys)
	if errCreate != nil {
		msg := "Create Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	return nil
}
