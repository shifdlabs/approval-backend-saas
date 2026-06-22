package documentnumbers

import (
	request "Microservice/data/request/DocumentNumbers" // Untuk DocumentNumbersRequest
	response "Microservice/data/response/DocumentNumbers"
	"Microservice/helper/enums"
	"Microservice/model"
	"fmt"
	"strconv"
	"time"

	// Untuk DocumentNumbersResponse

	// Untuk UserResponse
	"Microservice/helper"
	repository "Microservice/repository/DocumentNumbers" // Untuk DocumentNumbersRepository
	numberingFormatRepository "Microservice/repository/NumberingFormat"

	"strings"

	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
)

type DocumentNumbersServiceImpl struct {
	DocumentNumbersRepository repository.DocumentNumbersRepository
	NumberingFormatRepository numberingFormatRepository.NumberingFormatRepository
	Validate                  *validator.Validate
}

func NewDocumentNumbersServiceImpl(
	documentNumbersRepository repository.DocumentNumbersRepository,
	numberingFormatRepository numberingFormatRepository.NumberingFormatRepository,
	validate *validator.Validate) DocumentNumbersService {
	return &DocumentNumbersServiceImpl{
		DocumentNumbersRepository: documentNumbersRepository,
		NumberingFormatRepository: numberingFormatRepository,
		Validate:                  validate,
	}
}

func (t DocumentNumbersServiceImpl) Create(request request.DocumentNumbersRequest, userId string, document *model.Document, state enums.DocumentNumberState, orgID string) *helper.ErrorModel {

	// Get The Numbering Format
	numberingFormatData, fetchError := t.NumberingFormatRepository.Get(request.NumberingFormatID, orgID)
	if fetchError != nil {
		return fetchError
	}

	// If numbering format is nil (empty publicationValue), skip document number creation
	if numberingFormatData == nil {
		return nil
	}

	var groupId *string

	if *numberingFormatData.IncrementByGroup == true {
		strGroupID := numberingFormatData.GroupID.String()
		groupId = &strGroupID
	} else {
		groupId = nil
	}

	cancelledNumber, errCancelled := t.DocumentNumbersRepository.GetCancelled(numberingFormatData.ID.String(), groupId, orgID)
	if errCancelled != nil {
		return errCancelled
	}

	if cancelledNumber != nil {
		uuidStr, _ := uuid.FromString(request.NumberingFormatID)

		if cancelledNumber.ID == &uuidStr {
			number := t.UpdateMonthAndYear(numberingFormatData.Format, cancelledNumber.Value, numberingFormatData.Separator)
			cancelledNumber.Value = *number
			cancelledNumber.UserId = &userId
			cancelledNumber.State = int(state)
			cancelledNumber.Document = document
			err := t.DocumentNumbersRepository.Create(*cancelledNumber)
			if err != nil {
				return err
			}
		} else {
			number := t.GetNumberValue(numberingFormatData.Format, cancelledNumber.Value, numberingFormatData.Separator)
			newNumber := t.GenerateNewNumber(numberingFormatData.Format, *number, numberingFormatData.Separator)
			cancelledNumber.Value = *newNumber
			cancelledNumber.NumberingFormat = numberingFormatData
			cancelledNumber.UserId = &userId
			cancelledNumber.State = int(state)
			cancelledNumber.Document = document
			err := t.DocumentNumbersRepository.Create(*cancelledNumber)

			if err != nil {
				return err
			}
		}
	} else {
		countTotal, errCount := t.DocumentNumbersRepository.GetTotal(numberingFormatData.ID.String(), groupId, orgID)
		if errCount != nil {
			return errCount
		}

		currentNumber := strconv.Itoa(int(*countTotal) + 1)
		newNumber := t.GenerateNewNumber(numberingFormatData.Format, currentNumber, numberingFormatData.Separator)
		data := model.DocumentNumbers{
			UserId:          &userId,
			Value:           *newNumber,
			State:           int(state),
			NumberingFormat: numberingFormatData,
			Document:        document,
		}

		err := t.DocumentNumbersRepository.Create(data)

		if err != nil {
			return err
		}
	}

	return nil
}

func (t DocumentNumbersServiceImpl) Update(id string, document *model.Document, state enums.DocumentNumberState, orgID string) *helper.ErrorModel {
	res, err := t.DocumentNumbersRepository.Get(id, orgID)
	if err != nil {
		return err
	}

	res.State = int(state)
	res.Document = document
	res.DocumentID = (*uuid.UUID)(&document.ID)

	errUpdate := t.DocumentNumbersRepository.Update(*res, orgID)

	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

func (t DocumentNumbersServiceImpl) GetNumberValue(format string, value string, separator string) *string {
	formatArr := strings.Split(format, separator)
	indexOfNumber := indexOf(formatArr, string(enums.Number))

	if indexOfNumber != -1 {
		valueArr := strings.Split(value, separator)
		number := valueArr[indexOfNumber]
		return &number
	} else {
		return nil
	}
}

// Case: Generate New Number or Found cancelled number but with the different format
func (t DocumentNumbersServiceImpl) GenerateNewNumber(format string, number string, separator string) *string {
	formatArr := strings.Split(format, separator)
	value := make([]string, len(formatArr))
	for i, v := range formatArr {
		if v == string(enums.Year) {
			value[i] = GetCurrentYear()
		} else if v == string(enums.MonthNumber) {
			value[i] = GetCurrentMonth(enums.MonthNumber)
		} else if v == string(enums.MonthRoman) {
			value[i] = GetCurrentMonth(enums.MonthRoman)
		} else if v == string(enums.Number) {
			value[i] = number
		} else {
			value[i] = v
		}
	}

	result := strings.Join(value, separator)
	return &result
}

// Case: Found cancelled number but with the same format
func (t DocumentNumbersServiceImpl) UpdateMonthAndYear(format string, value string, separator string) *string {
	formatArr := strings.Split(format, separator)
	var indexOfYear int = -1
	var indexOfMonthNumber int = -1
	var indexOfMonthRoman int = -1

	if contains(formatArr, string(enums.Year)) {
		indexOfYear = indexOf(formatArr, string(enums.Year))
	}

	if contains(formatArr, string(enums.MonthNumber)) {
		indexOfMonthNumber = indexOf(formatArr, string(enums.MonthNumber))
	}

	if contains(formatArr, string(enums.MonthRoman)) {
		indexOfMonthRoman = indexOf(formatArr, string(enums.MonthRoman))
	}

	valueArr := strings.Split(value, separator)

	if indexOfYear != -1 {
		valueArr[indexOfYear] = GetCurrentYear()
	}

	if indexOfMonthNumber != -1 {
		valueArr[indexOfMonthNumber] = GetCurrentMonth(enums.MonthNumber)
	}

	if indexOfMonthRoman != -1 {
		valueArr[indexOfMonthRoman] = GetCurrentMonth(enums.MonthRoman)
	}

	result := strings.Join(valueArr, separator)

	return &result
}

func GetCurrentYear() string {
	currentYear := time.Now().Year()
	return strconv.Itoa(currentYear)
}

func GetCurrentMonth(format enums.FormatCode) string {
	romanMonths := []string{
		"I", "II", "III", "IV", "V", "VI",
		"VII", "VIII", "IX", "X", "XI", "XII",
	}

	now := time.Now()
	month := int(now.Month())

	switch format {
	case enums.MonthRoman:
		return romanMonths[month-1]
	case enums.MonthNumber:
		return fmt.Sprintf("%02d", month) // zero-padded
	default:
		return fmt.Sprintf("Invalid format: %s", format)
	}
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func (t DocumentNumbersServiceImpl) Get(id string, orgID string) (*response.DocumentNumbersResponse, *helper.ErrorModel) {
	result, fetchError := t.DocumentNumbersRepository.Get(id, orgID)
	if fetchError != nil {
		return nil, fetchError
	} else {
		response := t.convertDocumentNumbersToDocumentNumbersResponse(*result)
		return &response, nil
	}
}

func (t DocumentNumbersServiceImpl) GetByDocumentID(id uuid.UUID, orgID string) (*response.DocumentNumbersResponse, *helper.ErrorModel) {
	result, fetchError := t.DocumentNumbersRepository.GetByDocumentID(id, orgID)
	helper.PrintValue("Rezz A", "Rezz A")
	if fetchError != nil {
		helper.PrintValue("Rezz A1", "Rezz A")
		return nil, fetchError
	} else if result == nil {
		helper.PrintValue("Rezz A2", "Rezz A")
		return nil, nil
	} else {
		helper.PrintValue("Rezz A3", "Rezz A")
		response := t.convertDocumentNumbersToDocumentNumbersResponse(*result)
		return &response, nil
	}
}

func (t DocumentNumbersServiceImpl) GetAllByUserId(userId string, orgID string) ([]response.DocumentNumbersResponse, *helper.ErrorModel) {
	result, fetchError := t.DocumentNumbersRepository.GetAllByUserID(userId, orgID)
	if fetchError != nil {
		return nil, fetchError
	} else {
		return t.mapDocumentNumbersToDocumentNumbersResponse(result), nil
	}
}

func (t DocumentNumbersServiceImpl) GetAll(orgID string) ([]response.DocumentNumbersResponse, *helper.ErrorModel) {
	result, fetchError := t.DocumentNumbersRepository.GetAll(orgID)
	if fetchError != nil {
		return nil, fetchError
	} else {
		return t.mapDocumentNumbersToDocumentNumbersResponse(result), nil
	}
}

func (t DocumentNumbersServiceImpl) Delete(id string, orgID string) *helper.ErrorModel {
	errResponse := t.DocumentNumbersRepository.Delete(id, orgID)
	if errResponse != nil {
		return errResponse
	}

	return nil
}

func indexOf(slice []string, value string) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1 // not found
}
