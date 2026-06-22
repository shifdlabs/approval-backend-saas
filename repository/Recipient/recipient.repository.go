package recipient

import (
	"Microservice/helper"
	"Microservice/model"

	"gorm.io/gorm"
)

type RecipientRepository interface {
	Create(db gorm.DB, recipient []model.Recipient) *helper.ErrorModel
	GetRecipientsByDocId(id string) ([]model.Recipient, *helper.ErrorModel)
	GetAll() ([]model.Recipient, *helper.ErrorModel)
	Update(document model.Document, recipients []model.Recipient) error
}
