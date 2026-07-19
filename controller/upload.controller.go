package controller

import (
	"strings"
	"time"

	"Microservice/helper"
	"Microservice/pkg/s3presign"
	"Microservice/utils"

	"github.com/gin-gonic/gin"
)

// presignTTL is how long an issued S3 presigned URL stays valid. Kept short so a
// leaked URL is useless quickly; still ample for a single browser upload.
const presignTTL = 5 * time.Minute

// allowedUploadPrefixes restricts which object key prefixes the frontend may
// obtain a presigned URL for (AUDIT SEC-01, least-privilege). These match the
// folders the app actually writes to: attachments, signatures, and app assets
// (company logo / letter head).
var allowedUploadPrefixes = []string{
	"document-attachments/",
	"signatures/",
	"app-assets/",
}

// UploadController issues short-lived S3 presigned URLs so the browser can
// upload/delete objects directly without ever holding IAM credentials.
type UploadController struct {
	presigner *s3presign.Presigner
}

func NewUploadController(presigner *s3presign.Presigner) *UploadController {
	return &UploadController{presigner: presigner}
}

type presignRequest struct {
	Key         string `json:"key" validate:"required"`
	Method      string `json:"method"`      // "PUT" (default) or "DELETE"
	ContentType string `json:"contentType"` // advisory; not signed
}

type presignResponse struct {
	URL       string `json:"url"`
	Method    string `json:"method"`
	ExpiresIn int    `json:"expiresIn"` // seconds
}

// Presign validates the requested object key against the allowlist and returns a
// presigned URL. It requires the JWKS auth + subscription middleware upstream,
// so only authenticated org members can reach it.
func (controller *UploadController) Presign(ctx *gin.Context) {
	if !controller.presigner.Enabled() {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 503, Message: "File storage is not configured"})
		return
	}

	var payload presignRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure"})
		return
	}
	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	method := strings.ToUpper(strings.TrimSpace(payload.Method))
	if method == "" {
		method = "PUT"
	}
	if method != "PUT" && method != "DELETE" {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Unsupported method"})
		return
	}

	if !isAllowedKey(payload.Key) {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 403, Message: "Object key is not allowed"})
		return
	}

	url, err := controller.presigner.Presign(method, payload.Key, presignTTL)
	if err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 500, Message: "Failed to sign upload URL"})
		return
	}

	utils.SuccessResponse(ctx, presignResponse{
		URL:       url,
		Method:    method,
		ExpiresIn: int(presignTTL.Seconds()),
	})
}

// isAllowedKey rejects path traversal and keys outside the allowlisted prefixes.
func isAllowedKey(key string) bool {
	if key == "" || strings.HasPrefix(key, "/") || strings.Contains(key, "..") {
		return false
	}
	for _, prefix := range allowedUploadPrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}
