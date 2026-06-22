package authentication

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"time"

	"Microservice/config"
	model "Microservice/data/model/Authentication"
	authentication "Microservice/data/request/Authentication"
	"Microservice/helper"
	dbModel "Microservice/model"
	failedLoginAttemptRepository "Microservice/repository/FailedLoginAttempt"
	passwordResetTokenRepository "Microservice/repository/PasswordResetToken"
	userRepository "Microservice/repository/User"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type AuthServiceImpl struct {
	UserRepository               userRepository.UserRepository
	PasswordResetTokenRepository passwordResetTokenRepository.PasswordResetTokenRepository
	FailedLoginAttemptRepository failedLoginAttemptRepository.FailedLoginAttemptRepository
	Validate                     *validator.Validate
	Config                       config.Config
}

func NewAuthServiceImpl(
	userRepository userRepository.UserRepository,
	passwordResetTokenRepository passwordResetTokenRepository.PasswordResetTokenRepository,
	failedLoginAttemptRepository failedLoginAttemptRepository.FailedLoginAttemptRepository,
	validate *validator.Validate,
	cfg config.Config,
) AuthService {
	return &AuthServiceImpl{
		UserRepository:               userRepository,
		PasswordResetTokenRepository: passwordResetTokenRepository,
		FailedLoginAttemptRepository: failedLoginAttemptRepository,
		Validate:                     validate,
		Config:                       cfg,
	}
}

func generateSecureToken() (rawToken string, tokenHash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return
	}
	rawToken = hex.EncodeToString(b)
	h := sha256.Sum256([]byte(rawToken))
	tokenHash = hex.EncodeToString(h[:])
	return
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func isPasswordComplex(password string) bool {
	return regexp.MustCompile(`[A-Z]`).MatchString(password) &&
		regexp.MustCompile(`[a-z]`).MatchString(password) &&
		regexp.MustCompile(`[0-9]`).MatchString(password) &&
		regexp.MustCompile(`[!@#$%^&*()\-_=+\[\]{};:'",.<>/?\\|` + "`" + `~]`).MatchString(password)
}

func (t AuthServiceImpl) Login(payload authentication.LogInRequest) (model.LoginResult, *helper.ErrorModel) {
	// Step 1: Get user by email
	// NOTE: this method is unrouted in Phase 2 (SIS now owns login) but kept
	// compiling. No org_id exists pre-auth, hence the unscoped lookup.
	user, err := t.UserRepository.GetByEmailUnscoped(payload.Email)
	if err != nil {
		return model.LoginResult{
			AccessToken:  "",
			RefreshToken: "",
			User:         nil,
		}, err
	}

	// Step 2: Check Access field first
	if !user.Access {
		msg := "Your account access has been disabled. Please contact administrator."
		return model.LoginResult{
			AccessToken:  "",
			RefreshToken: "",
			User:         nil,
		}, &helper.ErrorModel{Code: 403, Message: msg}
	}

	// Step 3: Check if account is locked
	if user.IsLocked {
		msg := "Your account is locked due to multiple failed login attempts. Please contact administrator."
		return model.LoginResult{
			AccessToken:  "",
			RefreshToken: "",
			User:         nil,
		}, &helper.ErrorModel{Code: 403, Message: msg}
	}

	// Step 4: Verify password
	errVerifyPassword := helper.VerifyPassword(user.Password, payload.Password)
	if errVerifyPassword != nil {
		// Password is incorrect, record failed attempt
		now := time.Now()
		failedAttempt := dbModel.FailedLoginAttempt{
			UserID:      user.ID,
			AttemptedAt: &now,
		}

		// Save failed attempt to database
		errCreate := t.FailedLoginAttemptRepository.Create(failedAttempt)
		if errCreate != nil {
			// Log error but continue with login flow
			helper.GetFileAndLine(errCreate)
		}

		// Count total failed attempts for this user
		count, errCount := t.FailedLoginAttemptRepository.CountByUserId(user.ID.String())
		if errCount != nil {
			msg := "Incorrect password"
			return model.LoginResult{
				AccessToken:  "",
				RefreshToken: "",
				User:         nil,
			}, helper.ErrorCatcher(errVerifyPassword, 400, &msg)
		}

		// Check if attempts >= 3, lock the account
		if count >= 3 {
			// Lock the user account and disable access
			lockTime := time.Now()
			user.IsLocked = true
			user.LockTimestamp = &lockTime
			user.Access = false

			// Update user in database
			errUpdate := t.UserRepository.Update(*user, user.OrganizationID.String())
			if errUpdate != nil {
				// Log error but still return locked message
				helper.GetFileAndLine(errUpdate)
			}

			msg := "Too many failed login attempts. Your account has been locked. Please contact administrator."
			return model.LoginResult{
				AccessToken:  "",
				RefreshToken: "",
				User:         nil,
			}, helper.ErrorCatcher(errVerifyPassword, 403, &msg)
		}

		// Return error with remaining attempts
		remainingAttempts := 3 - count
		msg := strconv.Itoa(int(remainingAttempts)) + " attempt(s) remaining before account lock."
		return model.LoginResult{
			AccessToken:  "",
			RefreshToken: "",
			User:         nil,
		}, helper.ErrorCatcher(errVerifyPassword, 400, &msg)
	}

	// Step 5: Password is correct - Clear failed attempts and unlock if needed
	errDelete := t.FailedLoginAttemptRepository.DeleteByUserId(user.ID.String())
	if errDelete != nil {
		// Log error but continue with successful login
		helper.GetFileAndLine(errDelete)
	}

	// Reset lock status and access if user was locked
	if user.IsLocked {
		user.IsLocked = false
		user.LockTimestamp = nil
		user.Access = true
		errUpdate := t.UserRepository.Update(*user, user.OrganizationID.String())
		if errUpdate != nil {
			// Log error but continue with successful login
			helper.GetFileAndLine(errUpdate)
		}
	}

	// Step 6: Generate tokens and return success
	accessToken, _ := helper.GenerateAccessToken(user.ID.String())
	refreshToken, _ := helper.GenerateRefreshToken(user.ID.String())

	return model.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (t AuthServiceImpl) ForgotPassword(email string) *helper.ErrorModel {
	ctx := context.Background()
	rateLimitKey := "rate_limit:forgot_password:" + email
	count, redisErr := config.RedisClient.Incr(ctx, rateLimitKey).Result()
	if redisErr == nil {
		if count == 1 {
			config.RedisClient.Expire(ctx, rateLimitKey, time.Hour)
		}
		if count > 3 {
			return &helper.ErrorModel{Code: 429, Message: "Too many requests. Please try again later."}
		}
	}

	// No org_id in context — this route is unauthenticated by design.
	user, userErr := t.UserRepository.GetByEmailUnscoped(email)
	if userErr != nil {
		return nil
	}

	t.PasswordResetTokenRepository.InvalidateByUserID(user.ID.String())

	rawToken, tokenHash, genErr := generateSecureToken()
	if genErr != nil {
		msg := "Failed to generate reset token"
		return helper.ErrorCatcher(genErr, 500, &msg)
	}

	expiresAt := time.Now().Add(30 * time.Minute)
	resetToken := dbModel.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: &expiresAt,
	}
	if storeErr := t.PasswordResetTokenRepository.Create(resetToken); storeErr != nil {
		return storeErr
	}

	resetLink := t.Config.FrontendURL + "/reset-password?token=" + rawToken
	emailErr := helper.SendResetPasswordEmail(
		t.Config.SMTPHost,
		t.Config.SMTPPort,
		t.Config.SMTPUser,
		t.Config.SMTPPassword,
		t.Config.SMTPFrom,
		user.Email,
		resetLink,
	)
	if emailErr != nil {
		t.PasswordResetTokenRepository.InvalidateByUserID(user.ID.String())
		msg := "Failed to send reset email. Please try again."
		return helper.ErrorCatcher(emailErr, 500, &msg)
	}

	return nil
}

func (t AuthServiceImpl) ResetPassword(token, newPassword string) *helper.ErrorModel {
	if !isPasswordComplex(newPassword) {
		return &helper.ErrorModel{Code: 422, Message: "Password does not meet the required complexity."}
	}

	tokenHash := hashToken(token)
	resetToken, _ := t.PasswordResetTokenRepository.GetByTokenHash(tokenHash)
	if resetToken == nil {
		return &helper.ErrorModel{Code: 400, Message: "Reset token is invalid or has expired."}
	}

	if resetToken.UsedAt != nil {
		return &helper.ErrorModel{Code: 400, Message: "Reset token has already been used."}
	}

	if time.Now().After(*resetToken.ExpiresAt) {
		return &helper.ErrorModel{Code: 400, Message: "Reset token is invalid or has expired."}
	}

	hashedPassword, hashErr := dbModel.HashPasswordString(newPassword)
	if hashErr != nil {
		msg := "Failed to process new password"
		return helper.ErrorCatcher(hashErr, 500, &msg)
	}

	// No org_id in context — this route is unauthenticated by design.
	user, userErr := t.UserRepository.GetUnscoped(resetToken.UserID.String(), false)
	if userErr != nil {
		return userErr
	}

	user.Password = *hashedPassword
	if updateErr := t.UserRepository.Update(*user, user.OrganizationID.String()); updateErr != nil {
		return updateErr
	}

	t.PasswordResetTokenRepository.MarkUsed(tokenHash)

	return nil
}
