package middleware

import (
	"Microservice/helper"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminOnly enforces that the caller's org_role JWT claim is "owner" or "admin".
// Must be placed AFTER JWKSAuth, which stores org_role in the Gin context.
// The *gorm.DB parameter is retained for call-site compatibility but is unused.
func AdminOnly(_ *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := helper.GetOrgRole(ctx)
		if role != "owner" && role != "admin" {
			helper.ResponseError(ctx, helper.CustomError{Code: 403, Message: "Forbidden."})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
