package recipient

import (
	request "Microservice/data/request/Recipient"
	"Microservice/helper"
	"Microservice/model"
	documentRepository "Microservice/repository/Document"
	repository "Microservice/repository/Recipient"

	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type RecipientServiceImpl struct {
	RecipientRepository repository.RecipientRepository
	DocumentRepository  documentRepository.DocumentRepository
	Validate            *validator.Validate
	Db                  *gorm.DB
}

func NewRecipientServiceImpl(
	recipientRepository repository.RecipientRepository,
	documentRepsitory documentRepository.DocumentRepository,
	Db *gorm.DB,
	validate *validator.Validate) RecipientService {
	return &RecipientServiceImpl{
		RecipientRepository: recipientRepository,
		DocumentRepository:  documentRepsitory,
		Db:                  Db,
		Validate:            validate,
	}
}

func (t RecipientServiceImpl) Create(request request.RecipientRequest, orgID string) *helper.ErrorModel {
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

	var recipients []model.Recipient
	for _, userId := range request.UserIds {
		userUuid, errorParse := uuid.FromString(userId)
		if errorParse != nil {
			msg := "Structure Error"
			return helper.ErrorCatcher(errStructure, 500, &msg)
		}

		recipients = append(recipients, model.Recipient{
			Document: document,
			UserID:   userUuid,
		})
	}

	errCreate := t.RecipientRepository.Create(*trx, recipients)
	if errCreate != nil {
		trx.Rollback()
		msg := "Create Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	return nil
}

func (t RecipientServiceImpl) Update(request request.RecipientRequest, orgID string) *helper.ErrorModel {
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

	var recipients []model.Recipient
	for _, userId := range request.UserIds {
		userUuid, errorParse := uuid.FromString(userId)
		if errorParse != nil {
			msg := "Structure Error"
			return helper.ErrorCatcher(errStructure, 500, &msg)
		}

		recipients = append(recipients, model.Recipient{
			Document: document,
			UserID:   userUuid,
		})
	}

	errCreate := t.RecipientRepository.Update(*document, recipients)
	if errCreate != nil {
		msg := "Create Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	return nil
}
