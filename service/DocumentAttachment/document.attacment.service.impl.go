package documentAttachment

import (
	response "Microservice/data/response/DocumentAttachment"
	"Microservice/helper"
	repository "Microservice/repository/DocumentAttachment"

	"github.com/go-playground/validator/v10"
)

type DocumentAttachmentServiceImpl struct {
	DocumentAttachmentRepository repository.DocumentAttachmentRepository
	Validate                     *validator.Validate
}

func NewDocumentAttachmentServiceImpl(
	documentRepository repository.DocumentAttachmentRepository,
	validate *validator.Validate) DocumentAttachmentService {
	return &DocumentAttachmentServiceImpl{
		DocumentAttachmentRepository: documentRepository,
		Validate:                     validate,
	}
}

func (t DocumentAttachmentServiceImpl) Get(id string, orgID string) (*response.DocumentAttachmentResponse, *helper.ErrorModel) {
	document, fetchError := t.DocumentAttachmentRepository.Get(id, orgID)
	if fetchError != nil {
		return nil, fetchError
	}

	if document == nil {
		return nil, nil
	}

	documentResponse := t.convertDocumentAttachmentToDocumentAttachmentResponse(*document)

	return &documentResponse, fetchError
}

func (t DocumentAttachmentServiceImpl) GetAll(orgID string) ([]response.DocumentAttachmentResponse, *helper.ErrorModel) {
	result, fetchError := t.DocumentAttachmentRepository.GetAll(orgID)
	if fetchError != nil {
		return nil, fetchError
	} else {
		return t.mapDocumentAttachmentToDocumentAttachmentResponse(result), nil
	}
}

func (t DocumentAttachmentServiceImpl) Delete(id string, orgID string) *helper.ErrorModel {
	errResponse := t.DocumentAttachmentRepository.Delete(id, orgID)
	if errResponse != nil {
		return errResponse
	}

	return nil
}
