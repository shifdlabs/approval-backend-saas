package helper

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// Phase 2: the Approval Backend no longer issues or refreshes JWTs — all token
// minting lives in SIS. The local RSA signing keys and signer functions have
// been removed. The functions below remain only so the (now-unrouted) legacy
// login/refresh handlers still compile; they are disabled and never mint tokens.

const (
	RefreshTTL = 365 * 24 * time.Hour // referenced by the legacy refresh handler
)

// errLocalIssuanceDisabled is returned by the disabled local token signers.
var errLocalIssuanceDisabled = errors.New("local token issuance is disabled in Phase 2; tokens are issued by SIS")

// GenerateAccessToken is disabled — SIS issues access tokens.
func GenerateAccessToken(userID string) (string, error) {
	return "", errLocalIssuanceDisabled
}

// GenerateRefreshToken is disabled — SIS issues refresh tokens.
func GenerateRefreshToken(userID string) (string, error) {
	return "", errLocalIssuanceDisabled
}

// ValidateOrRefreshAccess is disabled — token refresh is handled by SIS.
func ValidateOrRefreshAccess(accessToken, refreshToken string) (newAccess string, newRefresh string, err error) {
	return "", "", errLocalIssuanceDisabled
}

// ─── Identity extraction (now sourced from the Gin context) ──────────────────
// The JWKS auth middleware validates the SIS JWT and stores the identity claims
// in the Gin context. These helpers read the user id from there so the existing
// controllers keep working without re-parsing the token on every call.

// GetUserId returns the authenticated user id from the request context.
func GetUserId(ctx *gin.Context) (*string, error) {
	id := GetUserID(ctx)
	if id == "" {
		return nil, errors.New("user id not found in request context")
	}
	return &id, nil
}

// GetUserUUID returns the authenticated user id parsed as a UUID, or nil.
func GetUserUUID(ctx *gin.Context) *uuid.UUID {
	id := GetUserID(ctx)
	if id == "" {
		return nil
	}
	userUUID, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse user id from context to uuid"
		ErrorLog(err, 500, &msg)
		return nil
	}
	return &userUUID
}
