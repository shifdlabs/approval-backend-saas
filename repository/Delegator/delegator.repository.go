package delegator

import (
	"Microservice/helper"
	"Microservice/model"
	"time"
)

type DelegatorRepository interface {
	Create(delegator model.Delegator, orgID string) *helper.ErrorModel
	GetAllByOwnerID(ownerID string, orgID string) ([]model.Delegator, *helper.ErrorModel)
	GetByID(id string, orgID string) (*model.Delegator, *helper.ErrorModel)
	Update(delegator model.Delegator, orgID string) *helper.ErrorModel
	Delete(id string, orgID string) *helper.ErrorModel
	GetActiveDelegationByOwnerID(ownerID string, date time.Time, orgID string) (*model.Delegator, *helper.ErrorModel)
	GetOwnerIDsByDelegateID(delegateID string, date time.Time, orgID string) ([]string, *helper.ErrorModel)
	HasOverlappingDelegation(ownerID string, startDate time.Time, endDate time.Time, excludeID *string, orgID string) (bool, *helper.ErrorModel)
}
