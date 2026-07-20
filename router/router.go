package router

import (
	"Microservice/config/middleware"
	"Microservice/controller"
	"Microservice/helper"
	"Microservice/pkg/jwks"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// parseCORSOrigins turns the comma-separated CORS_ALLOWED_ORIGINS env value into
// an allowlist set (AUDIT SEC-10). When empty it falls back to the built-in
// local-dev + production defaults so existing deployments keep working.
func parseCORSOrigins(raw string) map[string]bool {
	allowed := map[string]bool{}
	for _, o := range strings.Split(raw, ",") {
		if o = strings.TrimSpace(o); o != "" {
			allowed[o] = true
		}
	}
	if len(allowed) == 0 {
		allowed = map[string]bool{
			"http://localhost:5173":  true,
			"http://localhost:5174":  true,
			"http://localhost:5176":  true,
			"https://alphasoftn.com": true,
		}
	}
	return allowed
}

func CORS(allowedOrigins map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set(
				"Access-Control-Allow-Headers",
				"Content-Type, Authorization, X-Requested-With",
			)
			c.Writer.Header().Set(
				"Access-Control-Allow-Methods",
				"GET, POST, PUT, DELETE, OPTIONS",
			)
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func NewRouter(
	Db *gorm.DB,
	jwksClient *jwks.JWKSClient,
	userController *controller.UserController,
	documentController *controller.DocumentController,
	documentHistoryController *controller.DocumentHistoryController,
	documentAttachmentController *controller.DocumentAttachmentController,
	documentSequenceController *controller.DocumentSequenceController,
	positionController *controller.PositionController,
	userLogController *controller.UserLogController,
	appSettingsController *controller.AppSettingsController,
	recipientController *controller.RecipientController,
	bookmarkController *controller.BookmarkController,
	numberingGroupController *controller.NumberingGroupController,
	numberingFormatController *controller.NumberingFormatController,
	documentNumbersController *controller.DocumentNumbersController,
	signatureController *controller.SignatureController,
	delegatorController *controller.DelegatorController,
	verificationController *controller.VerificationController,
	letterTemplateController *controller.LetterTemplateController,
	uploadController *controller.UploadController,
	corsOrigins string,
) *gin.Engine {
	service := gin.Default()
	service.Use(CORS(parseCORSOrigins(corsOrigins)))

	// Shared Phase 2 auth pipeline: validate the SIS JWT (JWKS), then enforce
	// org scope + the shifd-approval subscription. Applied to every protected group.
	authMW := middleware.JWKSAuth(jwksClient, Db)
	subMW := middleware.SubscriptionCheck()

	// Public health check — no auth.
	service.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	service.GET("", func(context *gin.Context) {
		context.JSON(http.StatusOK, "Router has initialized")
	})

	service.NoRoute(func(c *gin.Context) {
		helper.ResponseError(c, helper.CustomError{
			Code:    404,
			Message: "Not Found.",
		})
	})

	router := service.Group("/api")

	// Phase 2: ALL authentication (login, register, refresh, logout, and
	// password recovery/activation) is owned by SIS. The Approval Backend
	// exposes no /auth routes — the frontend talks to SIS directly for those.

	protectedUserRouter := router.Group("/user")
	protectedUserRouter.Use(authMW, subMW)
	protectedUserRouter.POST("", middleware.AdminOnly(Db), userController.Create)
	protectedUserRouter.GET("/profile", userController.Get)
	protectedUserRouter.GET("/:id", userController.GetUserByID)
	protectedUserRouter.GET("", userController.GetAll)
	protectedUserRouter.GET("/except-current", userController.GetAllUserExceptCurrent)
	protectedUserRouter.PUT("", middleware.AdminOnly(Db), userController.Update)
	protectedUserRouter.DELETE("/:id", middleware.AdminOnly(Db), userController.Delete)
	protectedUserRouter.DELETE("/deletes", middleware.AdminOnly(Db), userController.MultipleDelete)
	protectedUserRouter.PUT("/role", middleware.AdminOnly(Db), userController.UpdateRole)
	protectedUserRouter.PUT("/password", userController.UpdatePassword)
	protectedUserRouter.PUT("/access", middleware.AdminOnly(Db), userController.UpdateAccess)
	protectedUserRouter.PUT("/unlock/:userId", middleware.AdminOnly(Db), userController.UnlockUser)
	protectedUserRouter.PUT("/biodata", userController.UpdateBiodata)
	protectedUserRouter.PUT("/email", userController.UpdateEmail)
	protectedUserRouter.POST("/import/preview", middleware.AdminOnly(Db), userController.PreviewImport)
	protectedUserRouter.POST("/import/bulk", middleware.AdminOnly(Db), userController.BulkImport)

	protectedDocumentRouter := router.Group("/document")
	protectedDocumentRouter.Use(authMW, subMW)
	protectedDocumentRouter.POST("", documentController.Create)
	protectedDocumentRouter.PUT("", documentController.Update)
	protectedDocumentRouter.GET("", documentController.GetAll)
	protectedDocumentRouter.GET("/references/:q", documentController.GetAllReferences)
	protectedDocumentRouter.GET("/:id", documentController.Get)
	protectedDocumentRouter.GET("/detail/:id", documentController.GetDetailPreview)
	protectedDocumentRouter.GET("/edit/:id", documentController.GetDetailForEdit)
	protectedDocumentRouter.POST("/authorize", documentController.Authorize)
	protectedDocumentRouter.GET("/authorization", documentController.GetAllAuthorization)
	protectedDocumentRouter.GET("/inprogress", documentController.GetAllInProgress)
	protectedDocumentRouter.GET("/inbox", documentController.GetAllInbox)
	protectedDocumentRouter.GET("/rejected", documentController.GetAllRejected)
	protectedDocumentRouter.GET("/dashboard", documentController.GetDashboardSummary)
	protectedDocumentRouter.GET("/dashboard/deadlines", documentController.GetDeadlines)
	protectedDocumentRouter.GET("/dashboard/activities", documentController.GetRecentActivities)
	protectedDocumentRouter.GET("/dashboard/recent", documentController.GetRecentDocuments)
	protectedDocumentRouter.GET("/search", documentController.Search)
	protectedDocumentRouter.POST("/:id/recall", documentController.Recall)

	protectedDocumentHistoryRouter := router.Group("/documenthistory")
	protectedDocumentHistoryRouter.Use(authMW, subMW)
	protectedDocumentHistoryRouter.GET("", documentHistoryController.GetAll)
	protectedDocumentHistoryRouter.GET("/:id", documentHistoryController.Get)
	protectedDocumentHistoryRouter.GET("/rejected", documentHistoryController.GetRejectedWithDocumentAndUser)

	protectedDocumentAttachmentRouter := router.Group("/documentattachment")
	protectedDocumentAttachmentRouter.Use(authMW, subMW)
	protectedDocumentAttachmentRouter.GET("", documentAttachmentController.GetAll)
	protectedDocumentAttachmentRouter.GET("/:id", documentAttachmentController.Get)
	protectedDocumentAttachmentRouter.DELETE("", documentAttachmentController.Delete)

	protectedDocumentRouter.GET("/complete", documentController.GetComplete)
	protectedDocumentRouter.GET("/draft", documentController.GetDraft)

	protectedUserLogRouter := router.Group("/userlogs")
	protectedUserLogRouter.Use(authMW, subMW)
	protectedUserLogRouter.GET("", userLogController.GetAll)
	protectedUserLogRouter.GET("/export", userLogController.Export)

	protectedDocumentSequenceRouter := router.Group("/documentsequence")
	protectedDocumentSequenceRouter.Use(authMW, subMW)
	//protectedDocumentSequenceRouter.GET("", documentSequenceController.GetAll)
	protectedDocumentSequenceRouter.GET("/:id", documentSequenceController.Get)
	protectedDocumentSequenceRouter.GET("/progress", documentSequenceController.GetProgress)

	protectedAppSettingsRouter := router.Group("/appsettings")
	protectedAppSettingsRouter.Use(authMW, subMW)
	protectedAppSettingsRouter.GET("", appSettingsController.GetAll)
	protectedAppSettingsRouter.PUT("", middleware.AdminOnly(Db), appSettingsController.Update)

	protectedPositionRouter := router.Group("/position")
	protectedPositionRouter.Use(authMW, subMW)
	protectedPositionRouter.GET("", positionController.GetAll)
	protectedPositionRouter.GET("/:id", positionController.Get)
	protectedPositionRouter.PUT("", middleware.AdminOnly(Db), positionController.Update)
	protectedPositionRouter.POST("", middleware.AdminOnly(Db), positionController.Create)
	protectedPositionRouter.DELETE("/:id", middleware.AdminOnly(Db), positionController.Delete)

	protectedBookmarkRouter := router.Group("/bookmark")
	protectedBookmarkRouter.Use(authMW, subMW)
	protectedBookmarkRouter.POST("/add", bookmarkController.AddBookmarkHandler)
	protectedBookmarkRouter.POST("/remove", bookmarkController.RemoveBookmarkHandler)
	protectedBookmarkRouter.POST("/status", bookmarkController.IsBookmarkedHandler)
	protectedBookmarkRouter.GET("/documents", bookmarkController.GetAllBookmarksWithDocumentsHandler)

	protectedNumberingGroupRouter := router.Group("/numbering/group")
	protectedNumberingGroupRouter.Use(authMW, subMW)
	protectedNumberingGroupRouter.GET("", numberingGroupController.GetAll)
	protectedNumberingGroupRouter.GET("/:id", numberingGroupController.Get)
	protectedNumberingGroupRouter.POST("", middleware.AdminOnly(Db), numberingGroupController.Create)
	protectedNumberingGroupRouter.DELETE("/:id", middleware.AdminOnly(Db), numberingGroupController.Delete)

	protectedNumberingFormatRouter := router.Group("/numbering/format")
	protectedNumberingFormatRouter.Use(authMW, subMW)
	protectedNumberingFormatRouter.GET("", numberingFormatController.GetAll)
	protectedNumberingFormatRouter.GET("/grouped", numberingFormatController.GetAllWithGrouped)
	protectedNumberingFormatRouter.POST("", middleware.AdminOnly(Db), numberingFormatController.Create)
	protectedNumberingFormatRouter.DELETE("/:id", middleware.AdminOnly(Db), numberingFormatController.Delete)

	protectedDocumentNumberRouter := router.Group("/document/number")
	protectedDocumentNumberRouter.Use(authMW, subMW)
	protectedDocumentNumberRouter.POST("", documentNumbersController.Create)
	protectedDocumentNumberRouter.GET("", documentNumbersController.GetAll)
	protectedDocumentNumberRouter.GET("/user", documentNumbersController.GetAllByUserId)
	protectedDocumentNumberRouter.DELETE("/:id", documentNumbersController.Delete)

	protectedSignatureRouter := router.Group("/signature")
	protectedSignatureRouter.Use(authMW, subMW)
	protectedSignatureRouter.GET("", signatureController.GetAll)
	protectedSignatureRouter.POST("", signatureController.Create)
	protectedSignatureRouter.PUT("/:userId", signatureController.Update)
	protectedSignatureRouter.DELETE("/:userId", signatureController.Delete)
	protectedSignatureRouter.GET("/:userId", signatureController.GetByUserId)

	protectedDelegatorRouter := router.Group("/delegator")
	protectedDelegatorRouter.Use(authMW, subMW)
	protectedDelegatorRouter.GET("", delegatorController.GetAll)
	protectedDelegatorRouter.POST("", delegatorController.Create)
	protectedDelegatorRouter.PUT("/:id", delegatorController.Update)
	protectedDelegatorRouter.DELETE("/:id", delegatorController.Delete)

	// S3 presigned-URL issuer (AUDIT SEC-01): authenticated org members exchange
	// an object key for a short-lived presigned PUT/DELETE URL, so the browser
	// never holds IAM credentials.
	protectedUploadRouter := router.Group("/upload")
	protectedUploadRouter.Use(authMW, subMW)
	protectedUploadRouter.POST("/presign", uploadController.Presign)

	// Public verification route — no auth middleware. Rate-limited per IP to
	// blunt scraping of documents by UUID (AUDIT SEC-10).
	verificationRouter := router.Group("/verification")
	verificationRouter.Use(middleware.RateLimit(30, time.Minute))
	verificationRouter.GET("/:id", verificationController.GetVerification)

	// Letter template routes — GET all/by-id: semua auth user; CUD: admin only
	templateRouter := router.Group("/template")
	templateRouter.Use(authMW, subMW)
	templateRouter.GET("", letterTemplateController.GetAll)
	templateRouter.GET("/:id", letterTemplateController.GetByID)
	templateRouter.POST("", middleware.AdminOnly(Db), letterTemplateController.Create)
	templateRouter.PUT("/:id", middleware.AdminOnly(Db), letterTemplateController.Update)
	templateRouter.DELETE("/:id", middleware.AdminOnly(Db), letterTemplateController.Delete)

	return service
}
