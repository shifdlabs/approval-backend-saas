package helper

import "github.com/gin-gonic/gin"

// Gin context keys populated by the JWKS auth middleware from the SIS JWT claims.
// Every tenant-scoped request carries these after authentication.
const (
	ContextKeyUserID   = "user_id"  // JWT "sub"      → users.id
	ContextKeyOrgID    = "org_id"   // JWT "org_id"   → tenant scope
	ContextKeyOrgRole  = "org_role" // JWT "org_role" → owner|admin|member
	ContextKeyEmail    = "email"    // JWT "email"
	ContextKeyName     = "name"     // JWT "name"
	ContextKeyProducts = "products" // JWT "products" → active subscriptions
)

// ctxString reads a string value from the Gin context, returning "" if absent
// or not a string.
func ctxString(c *gin.Context, key string) string {
	if v, ok := c.Get(key); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetUserID returns the authenticated user id (JWT sub claim), or "".
func GetUserID(c *gin.Context) string { return ctxString(c, ContextKeyUserID) }

// GetOrgID returns the organization id (JWT org_id claim), or "".
func GetOrgID(c *gin.Context) string { return ctxString(c, ContextKeyOrgID) }

// GetOrgRole returns the org role (JWT org_role claim): owner|admin|member, or "".
func GetOrgRole(c *gin.Context) string { return ctxString(c, ContextKeyOrgRole) }

// GetEmail returns the user email (JWT email claim), or "".
func GetEmail(c *gin.Context) string { return ctxString(c, ContextKeyEmail) }

// GetName returns the user display name (JWT name claim), or "".
func GetName(c *gin.Context) string { return ctxString(c, ContextKeyName) }

// GetProducts returns the active product subscriptions (JWT products claim).
func GetProducts(c *gin.Context) []string {
	if v, ok := c.Get(ContextKeyProducts); ok {
		if s, ok := v.([]string); ok {
			return s
		}
	}
	return nil
}

// RequireOrgID extracts org_id from the Gin context for use by handlers before
// calling into the service layer (org_id must flow handler → service →
// repository as an explicit parameter — repositories must never read Gin
// context directly, see CLAUDE.md Step 4).
//
// SubscriptionCheck middleware already rejects requests with no org_id before
// they reach any handler, so the empty case here is a defensive fallback (e.g.
// a route wired without that middleware), not the expected path.
func RequireOrgID(c *gin.Context) (string, bool) {
	orgID := GetOrgID(c)
	if orgID == "" {
		ResponseError(c, CustomError{Code: 403, Message: "Request is not scoped to an organization"})
		c.Abort()
		return "", false
	}
	return orgID, true
}
