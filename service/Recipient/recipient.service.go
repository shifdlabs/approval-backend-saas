package recipient

import (
	request "Microservice/data/request/Recipient"
	"Microservice/helper"
)

type RecipientService interface {
	Create(request request.RecipientRequest, orgID string) *helper.ErrorModel
	Update(request request.RecipientRequest, orgID string) *helper.ErrorModel
}
