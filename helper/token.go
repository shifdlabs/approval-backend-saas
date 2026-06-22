package helper

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// ExtractToken pulls a bearer token from the request: first from the "token"
// query parameter (used by some link-based flows), then from the Authorization
// header. Returns "" when absent.
//
// NOTE: local RSA key loading and JWT signature validation used to live here.
// In Phase 2 the Approval Backend validates SIS-issued JWTs via the JWKS auth
// middleware (see config/middleware/jwks_middleware.go), so those functions
// have been removed.
func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}

	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}
