package carbonCopy

import (
	"Microservice/helper"
	"Microservice/model"

	"gorm.io/gorm"
)

type CarbonCopyRepository interface {
	Create(db gorm.DB, carbonCopy []model.CarbonCopy) *helper.ErrorModel
	GetCarbonCopysByDocId(id string) ([]model.CarbonCopy, *helper.ErrorModel)
	GetAll() ([]model.CarbonCopy, *helper.ErrorModel)
	Update(document model.Document, carbonCopies []model.CarbonCopy) error
}
