package passwordresettoken

import (
	"Microservice/helper"
	"Microservice/model"
	"time"

	"gorm.io/gorm"
)

type PasswordResetTokenRepositoryImpl struct {
	Db *gorm.DB
}

func NewPasswordResetTokenRepositoryImpl(db *gorm.DB) PasswordResetTokenRepository {
	return &PasswordResetTokenRepositoryImpl{Db: db}
}

func (r *PasswordResetTokenRepositoryImpl) Create(token model.PasswordResetToken) *helper.ErrorModel {
	result := r.Db.Create(&token)
	if result.Error != nil {
		msg := "Failed to create password reset token"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return nil
}

func (r *PasswordResetTokenRepositoryImpl) GetByTokenHash(hash string) (*model.PasswordResetToken, *helper.ErrorModel) {
	var token model.PasswordResetToken
	result := r.Db.Where("token_hash = ? AND deleted_at IS NULL", hash).First(&token)
	if result.Error != nil {
		return nil, nil
	}
	return &token, nil
}

func (r *PasswordResetTokenRepositoryImpl) InvalidateByUserID(userID string) *helper.ErrorModel {
	result := r.Db.Unscoped().Where("user_id = ?", userID).Delete(&model.PasswordResetToken{})
	if result.Error != nil {
		msg := "Failed to invalidate existing tokens"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return nil
}

func (r *PasswordResetTokenRepositoryImpl) MarkUsed(tokenHash string) *helper.ErrorModel {
	now := time.Now()
	result := r.Db.Model(&model.PasswordResetToken{}).Where("token_hash = ?", tokenHash).Update("used_at", &now)
	if result.Error != nil {
		msg := "Failed to mark token as used"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return nil
}
