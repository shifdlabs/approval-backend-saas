package middleware

import (
	"Microservice/helper"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminOnly ensures the caller has role=99 (admin). Must be placed AFTER the
// JWKSAuth middleware, which populates the user id in the Gin context from the
// validated SIS JWT.
//
// TODO(Phase 2 / Step 4): scope this lookup by organization_id and prefer the
// JWT org_role claim ("owner"/"admin") over the local users.role column.
func AdminOnly(DB *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := helper.GetUserID(ctx)
		if userID == "" {
			helper.ResponseError(ctx, helper.CustomError{Code: 401, Message: "Unauthorize."})
			ctx.Abort()
			return
		}

		var role int
		result := DB.Raw("SELECT role FROM users WHERE id = ? AND deleted_at IS NULL", userID).Scan(&role)
		if result.Error != nil || result.RowsAffected == 0 {
			helper.ResponseError(ctx, helper.CustomError{Code: 401, Message: "Unauthorize."})
			ctx.Abort()
			return
		}

		if role != 99 {
			helper.ResponseError(ctx, helper.CustomError{Code: 403, Message: "Forbidden."})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
