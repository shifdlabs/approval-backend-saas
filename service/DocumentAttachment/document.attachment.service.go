package documentAttachment

import (
	response "Microservice/data/response/DocumentAttachment"
	"Microservice/helper"
)

type DocumentAttachmentService interface {
	Get(id string, orgID string) (*response.DocumentAttachmentResponse, *helper.ErrorModel)
	GetAll(orgID string) ([]response.DocumentAttachmentResponse, *helper.ErrorModel)
	Delete(id string, orgID string) *helper.ErrorModel
}
