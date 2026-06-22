package delegator

import (
	response "Microservice/data/response/Delegator"
	"Microservice/model"
	"time"
)

func (s *DelegatorServiceImpl) mapToDelegatorResponses(delegators []model.Delegator) []response.DelegatorResponse {
	result := make([]response.DelegatorResponse, len(delegators))
	for i, d := range delegators {
		result[i] = s.mapToDelegatorResponse(d)
	}
	return result
}

func (s *DelegatorServiceImpl) mapToDelegatorResponse(d model.Delegator) response.DelegatorResponse {
	now := time.Now()
	isActive := !now.Before(d.StartDate) && !now.After(d.EndDate)

	resp := response.DelegatorResponse{
		ID:        d.ID,
		OwnerID:   d.OwnerID,
		StartDate: d.StartDate,
		EndDate:   d.EndDate,
		IsActive:  isActive,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}

	if d.Owner != nil {
		resp.Owner = &response.DelegatorUserInfo{
			ID:        d.Owner.ID,
			FirstName: d.Owner.FirstName,
			LastName:  d.Owner.LastName,
			Email:     d.Owner.Email,
		}
	}

	if d.Delegate != nil {
		resp.Delegate = &response.DelegatorUserInfo{
			ID:        d.Delegate.ID,
			FirstName: d.Delegate.FirstName,
			LastName:  d.Delegate.LastName,
			Email:     d.Delegate.Email,
		}
	}

	return resp
}
