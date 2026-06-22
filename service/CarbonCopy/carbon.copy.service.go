package carbonCopy

import (
	request "Microservice/data/request/CarbonCopy"
	"Microservice/helper"
)

type CarbonCopyService interface {
	Create(request request.CarbonCopyRequest, orgID string) *helper.ErrorModel
	Update(request request.CarbonCopyRequest, orgID string) *helper.ErrorModel
}
