package delegator

import (
	request "Microservice/data/request/Delegator"
	response "Microservice/data/response/Delegator"
	"Microservice/helper"
	"time"
)

type DelegatorService interface {
	Create(ownerID string, req request.CreateDelegatorRequest, orgID string) *helper.ErrorModel
	GetAll(ownerID string, orgID string) ([]response.DelegatorResponse, *helper.ErrorModel)
	Update(id string, ownerID string, req request.UpdateDelegatorRequest, orgID string) *helper.ErrorModel
	Delete(id string, ownerID string, orgID string) *helper.ErrorModel
	ResolveDelegate(ownerID string, date time.Time, orgID string) (string, *helper.ErrorModel)
	GetOwnerChainForDelegate(delegateID string, date time.Time, orgID string) ([]string, *helper.ErrorModel)
}
