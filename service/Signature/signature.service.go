package signature

import (
	signatureRequest "Microservice/data/request/Signature"
	signatureResponse "Microservice/data/response/Signature"
	"Microservice/helper"
)

type SignatureService interface {
	Create(request signatureRequest.CreateSignatureRequest, orgID string) *helper.ErrorModel
	Update(userId string, request signatureRequest.UpdateSignatureRequest, orgID string) *helper.ErrorModel
	Delete(userId string, orgID string) *helper.ErrorModel
	GetAll(orgID string) ([]signatureResponse.SignatureResponse, *helper.ErrorModel)
	GetByUserId(userId string, orgID string) (*signatureResponse.SignatureResponse, *helper.ErrorModel)
}
