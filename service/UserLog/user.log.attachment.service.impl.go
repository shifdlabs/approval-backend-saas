package userlog

import (
	response "Microservice/data/response/UserLog"
	"Microservice/helper"
	"Microservice/model"
	repository "Microservice/repository/UserLog"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/xuri/excelize/v2"
)

type UserLogServiceImpl struct {
	UserLogRepository repository.UserLogRepository
	Validate          *validator.Validate
}

func NewUserLogServiceImpl(
	documentRepository repository.UserLogRepository,
	validate *validator.Validate) UserLogService {
	return &UserLogServiceImpl{
		UserLogRepository: documentRepository,
		Validate:          validate,
	}
}

func (t UserLogServiceImpl) GetAll(orgID string) ([]response.UserLogResponse, *helper.ErrorModel) {
	result, fetchError := t.UserLogRepository.GetAll(orgID)
	if fetchError != nil {
		return nil, fetchError
	} else {
		return t.mapUserLogToUserLogResponse(result), nil
	}
}

func (t UserLogServiceImpl) CreateLog(log model.UserLog, orgID string) {
	t.UserLogRepository.Create(log, orgID)
}

func (t UserLogServiceImpl) Export(orgID string) ([]byte, *helper.ErrorModel) {
	rows, err := t.UserLogRepository.GetAll(orgID)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Activity Log"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	headers := []string{"No", "User Name", "Action", "Module", "Log Date", "Detail"}
	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, row := range rows {
		rowNum := i + 2
		logDate := ""
		if row.LogDate != nil {
			logDate = row.LogDate.Format("2006-01-02 15:04:05")
		}
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), row.UserName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), row.Action)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), row.Module)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), logDate)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNum), string(row.Log))
	}

	buf, writeErr := f.WriteToBuffer()
	if writeErr != nil {
		msg := "Failed to write Excel file"
		return nil, helper.ErrorCatcher(writeErr, 500, &msg)
	}
	return buf.Bytes(), nil
}
