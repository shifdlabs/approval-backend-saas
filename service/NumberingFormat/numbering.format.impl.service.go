package numberingformat

import (
	request "Microservice/data/request/NumberingFormat" // Untuk NumberingFormatRequest
	response "Microservice/data/response/NumberingFormat"

	// Untuk NumberingFormatResponse
	"Microservice/model"

	// Untuk UserResponse
	"Microservice/helper"
	repository "Microservice/repository/NumberingFormat" // Untuk NumberingFormatRepository
	numberingGroupRepository "Microservice/repository/NumberingGroup"

	"github.com/go-playground/validator/v10"
)

type NumberingFormatServiceImpl struct {
	NumberingFormatRepository repository.NumberingFormatRepository
	NumberingGroupRepository  numberingGroupRepository.NumberingGroupRepository
	Validate                  *validator.Validate
}

func NewNumberingFormatServiceImpl(
	numberingFormatRepository repository.NumberingFormatRepository,
	numberingGroupRepository numberingGroupRepository.NumberingGroupRepository,
	validate *validator.Validate) NumberingFormatService {
	return &NumberingFormatServiceImpl{
		NumberingFormatRepository: numberingFormatRepository,
		NumberingGroupRepository:  numberingGroupRepository,
		Validate:                  validate,
	}
}

func (t NumberingFormatServiceImpl) Create(request request.NumberingFormatRequest, orgID string) *helper.ErrorModel {
	group, errorResponse := t.NumberingGroupRepository.Get(request.GroupID, orgID)

	if errorResponse != nil {
		return errorResponse
	}

	data := model.NumberingFormat{
		Name:             request.Name,
		Format:           request.Format,
		Separator:        request.Separator,
		Group:            *group,
		IncrementByGroup: request.IncrementByGroup,
	}

	fetchError := t.NumberingFormatRepository.Create(data)
	if fetchError != nil {
		return fetchError
	}

	return nil
}

func (t NumberingFormatServiceImpl) GetAll(orgID string) ([]response.NumberingFormatResponse, *helper.ErrorModel) {
	result, fetchError := t.NumberingFormatRepository.GetAll(orgID)
	groups, fetchGroupError := t.NumberingGroupRepository.GetAll(orgID)

	if fetchError != nil || fetchGroupError != nil {
		if fetchError != nil {
			return nil, fetchError
		} else {
			return nil, fetchGroupError
		}
	} else {
		usedGroupNames := make(map[string]bool)
		numberingFormatsResult := make([]response.NumberingFormatResponse, len(result))
		for i, numberingFormat := range result {
			numberingFormatsResult[i] = response.NumberingFormatResponse{
				Id:                 numberingFormat.ID,
				Name:               numberingFormat.Name,
				Format:             numberingFormat.Format,
				IncrementedByGroup: *numberingFormat.IncrementByGroup,
				Separator:          numberingFormat.Separator,
				CreatedAt:          numberingFormat.CreatedAt,
				Group:              numberingFormat.Group.Name,
			}

			usedGroupNames[numberingFormat.Group.Name] = true
		}

		for _, group := range groups {
			if _, exists := usedGroupNames[group.Name]; !exists {
				numberingFormatsResult = append(numberingFormatsResult, response.NumberingFormatResponse{
					Id:                 nil,
					Name:               "No Data Available",
					Format:             "",
					IncrementedByGroup: false,
					Separator:          "",
					CreatedAt:          nil,
					Group:              group.Name,
				})
			}
		}

		return numberingFormatsResult, nil
	}
}

func (t NumberingFormatServiceImpl) GetAllWithGrouped(orgID string) ([]response.NumberingFormatByGroupResponse, *helper.ErrorModel) {
	result, fetchError := t.NumberingFormatRepository.GetAll(orgID)
	if fetchError != nil {
		return nil, fetchError
	}

	groupMap := make(map[string][]response.Format)
	for _, nf := range result {
		format := response.Format{
			Id:   nf.ID.String(),
			Name: nf.Name,
		}
		groupName := nf.Group.Name
		groupMap[groupName] = append(groupMap[groupName], format)
	}

	var mappedValue []response.NumberingFormatByGroupResponse
	for groupName, formats := range groupMap {
		mappedValue = append(mappedValue, response.NumberingFormatByGroupResponse{
			Group:   groupName,
			Formats: formats,
		})
	}

	return mappedValue, nil
}

func (t NumberingFormatServiceImpl) Delete(id string, orgID string) *helper.ErrorModel {
	errResponse := t.NumberingFormatRepository.Delete(id, orgID)
	if errResponse != nil {
		return errResponse
	}

	return nil
}
