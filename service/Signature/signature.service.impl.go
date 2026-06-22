package signature

import (
	signatureRequest "Microservice/data/request/Signature"
	signatureResponse "Microservice/data/response/Signature"
	"Microservice/helper"
	"Microservice/model"
	signatureRepository "Microservice/repository/Signature"
	"errors"

	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
)

type SignatureServiceImpl struct {
	SignatureRepository signatureRepository.SignatureRepository
	Validate            *validator.Validate
}

func NewSignatureServiceImpl(signatureRepository signatureRepository.SignatureRepository, validate *validator.Validate) SignatureService {
	return &SignatureServiceImpl{
		SignatureRepository: signatureRepository,
		Validate:            validate,
	}
}

func (t SignatureServiceImpl) Create(request signatureRequest.CreateSignatureRequest, orgID string) *helper.ErrorModel {
	err := t.Validate.Struct(request)
	if err != nil {
		msg := "Structure Error"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	// Check if signature already exists for this user
	existingSignature, _ := t.SignatureRepository.GetByUserId(request.UserID, orgID)
	if existingSignature != nil {
		msg := "Signature already exists for this user"
		return helper.ErrorCatcher(errors.New("duplicate signature"), 400, &msg)
	}

	userIdUUID, parseErr := uuid.FromString(request.UserID)
	if parseErr != nil {
		msg := "Failed to parse user ID"
		return helper.ErrorCatcher(parseErr, 500, &msg)
	}

	signature := &model.Signature{
		UserID:   &userIdUUID,
		ImageURL: request.ImageURL,
	}

	createErr := t.SignatureRepository.Create(signature, orgID)
	if createErr != nil {
		return createErr
	}

	return nil
}

func (t SignatureServiceImpl) Update(userId string, request signatureRequest.UpdateSignatureRequest, orgID string) *helper.ErrorModel {
	err := t.Validate.Struct(request)
	if err != nil {
		msg := "Structure Error"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	signature, getErr := t.SignatureRepository.GetByUserId(userId, orgID)
	if getErr != nil {
		return getErr
	}

	if signature == nil {
		msg := "Signature not found"
		return helper.ErrorCatcher(nil, 404, &msg)
	}

	signature.ImageURL = request.ImageURL

	updateErr := t.SignatureRepository.Update(signature)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (t SignatureServiceImpl) Delete(userId string, orgID string) *helper.ErrorModel {
	signature, getErr := t.SignatureRepository.GetByUserId(userId, orgID)
	if getErr != nil {
		return getErr
	}

	if signature == nil {
		msg := "Signature not found"
		return helper.ErrorCatcher(nil, 404, &msg)
	}

	deleteErr := t.SignatureRepository.Delete(signature.ID.String())
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func (t SignatureServiceImpl) GetAll(orgID string) ([]signatureResponse.SignatureResponse, *helper.ErrorModel) {
	signatures, getErr := t.SignatureRepository.GetAll(orgID)
	if getErr != nil {
		return nil, getErr
	}

	var responses []signatureResponse.SignatureResponse
	for _, signature := range signatures {
		response := signatureResponse.SignatureResponse{
			ID:       signature.ID.String(),
			UserID:   signature.UserID.String(),
			ImageURL: signature.ImageURL,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (t SignatureServiceImpl) GetByUserId(userId string, orgID string) (*signatureResponse.SignatureResponse, *helper.ErrorModel) {
	signature, getErr := t.SignatureRepository.GetByUserId(userId, orgID)
	if getErr != nil {
		return nil, getErr
	}

	if signature == nil {
		return nil, nil
	}

	response := &signatureResponse.SignatureResponse{
		ID:       signature.ID.String(),
		UserID:   signature.UserID.String(),
		ImageURL: signature.ImageURL,
	}

	return response, nil
}
