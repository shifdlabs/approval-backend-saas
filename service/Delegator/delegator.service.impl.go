package delegator

import (
	request "Microservice/data/request/Delegator"
	response "Microservice/data/response/Delegator"
	"Microservice/helper"
	"Microservice/model"
	repository "Microservice/repository/Delegator"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
)

const maxDelegationDepth = 10

type DelegatorServiceImpl struct {
	DelegatorRepository repository.DelegatorRepository
	Validate            *validator.Validate
}

func NewDelegatorServiceImpl(repo repository.DelegatorRepository, validate *validator.Validate) DelegatorService {
	return &DelegatorServiceImpl{
		DelegatorRepository: repo,
		Validate:            validate,
	}
}

func (s *DelegatorServiceImpl) Create(ownerID string, req request.CreateDelegatorRequest, orgID string) *helper.ErrorModel {
	if errs := helper.ValidateStruct(req); len(errs) > 0 {
		msg := "Invalid request data"
		return helper.ErrorCatcher(fmt.Errorf("validation failed"), 400, &msg)
	}

	startDate, errParse := time.Parse("2006-01-02", req.StartDate)
	if errParse != nil {
		msg := "Invalid start_date format. Use YYYY-MM-DD"
		return helper.ErrorCatcher(errParse, 400, &msg)
	}

	endDate, errParse := time.Parse("2006-01-02", req.EndDate)
	if errParse != nil {
		msg := "Invalid end_date format. Use YYYY-MM-DD"
		return helper.ErrorCatcher(errParse, 400, &msg)
	}

	if !endDate.After(startDate) {
		msg := "end_date must be after start_date"
		return helper.ErrorCatcher(fmt.Errorf("invalid date range"), 400, &msg)
	}

	ownerUUID, errParse := uuid.FromString(ownerID)
	if errParse != nil {
		msg := "Invalid owner ID"
		return helper.ErrorCatcher(errParse, 400, &msg)
	}

	delegateUUID, errParse := uuid.FromString(req.DelegateID)
	if errParse != nil {
		msg := "Invalid delegate_id"
		return helper.ErrorCatcher(errParse, 400, &msg)
	}

	if ownerUUID == delegateUUID {
		msg := "You cannot delegate to yourself"
		return helper.ErrorCatcher(fmt.Errorf("self delegation"), 400, &msg)
	}

	overlapping, errCheck := s.DelegatorRepository.HasOverlappingDelegation(ownerID, startDate, endDate, nil, orgID)
	if errCheck != nil {
		return errCheck
	}
	if overlapping {
		msg := "A delegation already exists for the specified date range. Only one delegate is allowed per date range."
		return helper.ErrorCatcher(fmt.Errorf("overlapping delegation"), 409, &msg)
	}

	delegator := model.Delegator{
		OwnerID:    ownerUUID,
		DelegateID: delegateUUID,
		StartDate:  startDate,
		EndDate:    endDate,
	}

	return s.DelegatorRepository.Create(delegator, orgID)
}

func (s *DelegatorServiceImpl) GetAll(ownerID string, orgID string) ([]response.DelegatorResponse, *helper.ErrorModel) {
	delegators, err := s.DelegatorRepository.GetAllByOwnerID(ownerID, orgID)
	if err != nil {
		return nil, err
	}
	return s.mapToDelegatorResponses(delegators), nil
}

func (s *DelegatorServiceImpl) Update(id string, ownerID string, req request.UpdateDelegatorRequest, orgID string) *helper.ErrorModel {
	if errs := helper.ValidateStruct(req); len(errs) > 0 {
		msg := "Invalid request data"
		return helper.ErrorCatcher(fmt.Errorf("validation failed"), 400, &msg)
	}

	existing, errGet := s.DelegatorRepository.GetByID(id, orgID)
	if errGet != nil {
		return errGet
	}

	if existing.OwnerID.String() != ownerID {
		msg := "You are not authorized to update this delegation"
		return helper.ErrorCatcher(fmt.Errorf("unauthorized"), 403, &msg)
	}

	startDate, errParse := time.Parse("2006-01-02", req.StartDate)
	if errParse != nil {
		msg := "Invalid start_date format. Use YYYY-MM-DD"
		return helper.ErrorCatcher(errParse, 400, &msg)
	}

	endDate, errParse := time.Parse("2006-01-02", req.EndDate)
	if errParse != nil {
		msg := "Invalid end_date format. Use YYYY-MM-DD"
		return helper.ErrorCatcher(errParse, 400, &msg)
	}

	if !endDate.After(startDate) {
		msg := "end_date must be after start_date"
		return helper.ErrorCatcher(fmt.Errorf("invalid date range"), 400, &msg)
	}

	ownerUUID, _ := uuid.FromString(ownerID)
	delegateUUID, errParse := uuid.FromString(req.DelegateID)
	if errParse != nil {
		msg := "Invalid delegate_id"
		return helper.ErrorCatcher(errParse, 400, &msg)
	}

	if ownerUUID == delegateUUID {
		msg := "You cannot delegate to yourself"
		return helper.ErrorCatcher(fmt.Errorf("self delegation"), 400, &msg)
	}

	idStr := id
	overlapping, errCheck := s.DelegatorRepository.HasOverlappingDelegation(ownerID, startDate, endDate, &idStr, orgID)
	if errCheck != nil {
		return errCheck
	}
	if overlapping {
		msg := "A delegation already exists for the specified date range. Only one delegate is allowed per date range."
		return helper.ErrorCatcher(fmt.Errorf("overlapping delegation"), 409, &msg)
	}

	existing.DelegateID = delegateUUID
	existing.StartDate = startDate
	existing.EndDate = endDate

	return s.DelegatorRepository.Update(*existing, orgID)
}

func (s *DelegatorServiceImpl) Delete(id string, ownerID string, orgID string) *helper.ErrorModel {
	existing, errGet := s.DelegatorRepository.GetByID(id, orgID)
	if errGet != nil {
		return errGet
	}

	if existing.OwnerID.String() != ownerID {
		msg := "You are not authorized to delete this delegation"
		return helper.ErrorCatcher(fmt.Errorf("unauthorized"), 403, &msg)
	}

	return s.DelegatorRepository.Delete(id, orgID)
}

// ResolveDelegate follows the delegation chain from ownerID on the given date
// and returns the ID of the final delegate who should act on their behalf.
// If no active delegation exists, returns ownerID itself.
func (s *DelegatorServiceImpl) ResolveDelegate(ownerID string, date time.Time, orgID string) (string, *helper.ErrorModel) {
	current := ownerID
	for depth := 0; depth < maxDelegationDepth; depth++ {
		delegation, err := s.DelegatorRepository.GetActiveDelegationByOwnerID(current, date, orgID)
		if err != nil {
			return "", err
		}
		if delegation == nil {
			return current, nil
		}
		next := delegation.DelegateID.String()
		if next == ownerID {
			// Circular delegation detected — return current to avoid infinite loop
			return current, nil
		}
		current = next
	}
	return current, nil
}

// GetOwnerChainForDelegate returns all ownerIDs whose delegation chain eventually resolves
// to delegateID on the given date. This is used to determine which documents a delegate can act on.
func (s *DelegatorServiceImpl) GetOwnerChainForDelegate(delegateID string, date time.Time, orgID string) ([]string, *helper.ErrorModel) {
	visited := map[string]bool{delegateID: true}
	result := []string{}
	queue := []string{delegateID}

	for len(queue) > 0 && len(result) < maxDelegationDepth*10 {
		current := queue[0]
		queue = queue[1:]

		ownerIDs, err := s.DelegatorRepository.GetOwnerIDsByDelegateID(current, date, orgID)
		if err != nil {
			return nil, err
		}

		for _, ownerID := range ownerIDs {
			if !visited[ownerID] {
				visited[ownerID] = true
				result = append(result, ownerID)
				queue = append(queue, ownerID)
			}
		}
	}

	return result, nil
}
