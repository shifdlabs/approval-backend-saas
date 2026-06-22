package lettertemplate

import (
	"Microservice/helper"
	"Microservice/model"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type LetterTemplateRepositoryImpl struct {
	Db *gorm.DB
}

func NewLetterTemplateRepositoryImpl(db *gorm.DB) LetterTemplateRepository {
	return &LetterTemplateRepositoryImpl{Db: db}
}

func (r *LetterTemplateRepositoryImpl) GetAll(orgID string) ([]model.LetterTemplate, *helper.ErrorModel) {
	var templates []model.LetterTemplate
	if err := r.Db.Where("organization_id = ?", orgID).Order("created_at DESC").Find(&templates).Error; err != nil {
		msg := "Failed to get letter templates"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}
	return templates, nil
}

func (r *LetterTemplateRepositoryImpl) GetByID(id string, orgID string) (*model.LetterTemplate, *helper.ErrorModel) {
	var template model.LetterTemplate
	if err := r.Db.Where("organization_id = ? AND id = ?", orgID, id).First(&template).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			msg := "Letter template not found"
			return nil, helper.ErrorCatcher(err, 404, &msg)
		}
		msg := "Failed to get letter template"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}
	return &template, nil
}

func (r *LetterTemplateRepositoryImpl) Create(template model.LetterTemplate) (*model.LetterTemplate, *helper.ErrorModel) {
	template.ID = uuid.NewV4()
	if err := r.Db.Create(&template).Error; err != nil {
		msg := "Failed to create letter template"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}
	return &template, nil
}

func (r *LetterTemplateRepositoryImpl) Update(id string, template model.LetterTemplate, orgID string) (*model.LetterTemplate, *helper.ErrorModel) {
	existing, errModel := r.GetByID(id, orgID)
	if errModel != nil {
		return nil, errModel
	}
	existing.Name = template.Name
	existing.Description = template.Description
	existing.Body = template.Body
	if err := r.Db.Save(existing).Error; err != nil {
		msg := "Failed to update letter template"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}
	return existing, nil
}

func (r *LetterTemplateRepositoryImpl) Delete(id string, orgID string) *helper.ErrorModel {
	existing, errModel := r.GetByID(id, orgID)
	if errModel != nil {
		return errModel
	}
	if err := r.Db.Delete(existing).Error; err != nil {
		msg := "Failed to delete letter template"
		return helper.ErrorCatcher(err, 500, &msg)
	}
	return nil
}
