package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Microservice/config"
	"Microservice/controller"
	"Microservice/model"
	"Microservice/pkg/jwks"
	"Microservice/pkg/s3presign"
	"Microservice/router"

	appSettingsRepository "Microservice/repository/AppSettings"
	bookmarkRepository "Microservice/repository/Bookmark"
	carbonCopiesRepository "Microservice/repository/CarbonCopy"
	delegatorRepository "Microservice/repository/Delegator"
	documentRepository "Microservice/repository/Document"
	documentAttachmentRepository "Microservice/repository/DocumentAttachment"
	documentHistoryRepository "Microservice/repository/DocumentHistory"
	documentNumbersRepository "Microservice/repository/DocumentNumbers"
	documentReferenceRepository "Microservice/repository/DocumentReference"
	documentSequenceRepository "Microservice/repository/DocumentSequence"
	failedLoginAttemptRepository "Microservice/repository/FailedLoginAttempt"
	numberingFormatRepository "Microservice/repository/NumberingFormat"
	numberingGroupRepository "Microservice/repository/NumberingGroup"
	positionRepository "Microservice/repository/Position"
	recipientRepository "Microservice/repository/Recipient"
	signatureRepository "Microservice/repository/Signature"
	userRepository "Microservice/repository/User"
	userLogRepository "Microservice/repository/UserLog"

	appSettingService "Microservice/service/AppSettings"
	bookmarkService "Microservice/service/Bookmark"
	delegatorService "Microservice/service/Delegator"
	slaService "Microservice/service/SLA"
	documentService "Microservice/service/Document"
	documentAttachmentService "Microservice/service/DocumentAttachment"
	documentHistoryService "Microservice/service/DocumentHistory"
	documentNumbersService "Microservice/service/DocumentNumbers"
	documentSequenceService "Microservice/service/DocumentSequence"
	numberingFormatService "Microservice/service/NumberingFormat"
	numberingGroupService "Microservice/service/NumberingGroup"
	positionService "Microservice/service/Position"
	recipientService "Microservice/service/Recipient"
	signatureService "Microservice/service/Signature"
	userService "Microservice/service/User"
	userLogService "Microservice/service/UserLog"
	letterTemplateRepository "Microservice/repository/LetterTemplate"
	letterTemplateService "Microservice/service/LetterTemplate"
	emailService "Microservice/service/Email"

	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func main() {
	envConf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	// Redis
	config.ConnectRedis(&envConf)

	// Database
	db := config.DatabaseConnection(&envConf)
	validate := validator.New()

	// SIS JWKS: fetch the SIS public key on startup (panics if unreachable) and
	// refresh it in the background to handle key rotation.
	jwksClient, _ := jwks.NewJWKSClient(envConf.SISJWKSURL, envConf.SISIssuer)
	jwksRefresh := envConf.JWKSRefreshInterval
	if jwksRefresh <= 0 {
		jwksRefresh = 24 * time.Hour
	}
	jwksClient.StartAutoRefresh(jwksRefresh)

	println("Message: Migrating Table... ")
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Document{})
	db.AutoMigrate(&model.Position{})
	db.AutoMigrate(&model.DocumentHistory{})
	db.AutoMigrate(&model.DocumentSequence{})
	db.AutoMigrate(&model.DocumentAttachment{})
	db.AutoMigrate(&model.Recipient{})
	db.AutoMigrate(&model.AppSettings{})
	db.AutoMigrate(&model.UserLog{})
	db.AutoMigrate(&model.CarbonCopy{})
	db.AutoMigrate(&model.NumberingGroup{})
	db.AutoMigrate(&model.NumberingFormat{})
	db.AutoMigrate(&model.DocumentNumbers{})
	db.AutoMigrate(&model.DocumentReference{})
	db.AutoMigrate(&model.Signature{})
	db.AutoMigrate(&model.NumberingGroup{}, &model.NumberingFormat{})
	db.AutoMigrate(&model.NumberingFormat{}, &model.DocumentNumbers{})
	db.AutoMigrate(&model.User{}, &model.Position{})
	db.AutoMigrate(&model.Document{}, &model.DocumentHistory{})
	db.AutoMigrate(&model.Document{}, &model.DocumentSequence{})
	db.AutoMigrate(&model.Document{}, &model.DocumentAttachment{})
	db.AutoMigrate(&model.Document{}, &model.Recipient{})
	db.AutoMigrate(&model.Document{}, &model.CarbonCopy{})
	db.AutoMigrate(&model.Document{}, &model.Bookmark{})
	db.AutoMigrate(&model.Document{}, &model.DocumentNumbers{})
	db.AutoMigrate(&model.User{}, &model.Signature{})
	db.AutoMigrate(&model.Delegator{})
	db.AutoMigrate(&model.User{}, &model.Delegator{})
	db.AutoMigrate(&model.LetterTemplate{})
	db.AutoMigrate(&model.Document{}, &model.LetterTemplate{})

	// Repositories
	userRepository := userRepository.NewUserRepositoryImpl(db)
	documentRepository := documentRepository.NewDocumentRepositoryImpl(db)
	documentHistoryRepository := documentHistoryRepository.NewDocumentHistoryRepositoryImpl(db)
	documentSequenceRepository := documentSequenceRepository.NewDocumentSequenceRepositoryImpl(db)
	documentAttachmentRepository := documentAttachmentRepository.NewDocumentAttachmentRepositoryImpl(db)
	positionRepositoy := positionRepository.NewPositionRepositoryImpl(db)
	userLogRepository := userLogRepository.NewUserLogRepositoryImpl(db)
	appSettingsRepository := appSettingsRepository.NewAppSettingsRepositoryImpl(db)
	recipientRepository := recipientRepository.NewRecipientRepositoryImpl(db)
	carbonCopiesRepository := carbonCopiesRepository.NewCarbonCopyRepositoryImpl(db)
	bookmarkRepository := bookmarkRepository.NewBookmarkRepositoryImpl(db)
	numberingGroupRepository := numberingGroupRepository.NewNumberingGroupRepositoryImpl(db)
	numberingFormatRepository := numberingFormatRepository.NewNumberingFormatRepositoryImpl(db)
	documentNumbersRepository := documentNumbersRepository.NewDocumentNumbersRepositoryImpl(db)
	documentReferenceRepository := documentReferenceRepository.NewDocumentReferenceRepositoryImpl(db)
	signatureRepository := signatureRepository.NewSignatureRepositoryImpl(db)
	failedLoginAttemptRepository := failedLoginAttemptRepository.NewFailedLoginAttemptRepositoryImpl(db)
	delegatorRepository := delegatorRepository.NewDelegatorRepositoryImpl(db)

	// Servic
	userService := userService.NewUserServiceImpl(userRepository, positionRepositoy, failedLoginAttemptRepository, validate)
	userLogService := userLogService.NewUserLogServiceImpl(userLogRepository, validate)
	documentSequenceService := documentSequenceService.NewDocumentSequenceServiceImpl(documentSequenceRepository, validate)
	emailSvc := emailService.NewEmailService(envConf.ResendAPIKey, envConf.EmailFrom, envConf.FrontendURL)
	documentService := documentService.NewDocumentServiceImpl(documentRepository, userRepository, documentSequenceRepository, documentAttachmentRepository, documentHistoryRepository, recipientRepository, carbonCopiesRepository, userLogRepository, documentNumbersRepository, documentReferenceRepository, signatureRepository, delegatorRepository, appSettingsRepository, emailSvc, envConf.FrontendURL, db, validate)
	documentHistoryService := documentHistoryService.NewDocumentHistoryServiceImpl(documentHistoryRepository, validate)
	documentAttachmentService := documentAttachmentService.NewDocumentAttachmentServiceImpl(documentAttachmentRepository, validate)
	positionService := positionService.NewPositionServiceImpl(positionRepositoy, validate)
	appSettingsService := appSettingService.NewAppSettingsServiceImpl(appSettingsRepository, validate)
	recipientService := recipientService.NewRecipientServiceImpl(recipientRepository, documentRepository, db, validate)
	bookmarkService := bookmarkService.NewBookmarkServiceImpl(bookmarkRepository, validate)
	numberingGroupService := numberingGroupService.NewNumberingGroupServiceImpl(numberingGroupRepository, validate)
	numberingFormatService := numberingFormatService.NewNumberingFormatServiceImpl(numberingFormatRepository, numberingGroupRepository, validate)
	documentNumbersService := documentNumbersService.NewDocumentNumbersServiceImpl(documentNumbersRepository, numberingFormatRepository, validate)
	signatureService := signatureService.NewSignatureServiceImpl(signatureRepository, validate)
	delegatorSvc := delegatorService.NewDelegatorServiceImpl(delegatorRepository, validate)
	slaSvc := slaService.NewSLAServiceImpl(appSettingsRepository, documentRepository, documentHistoryRepository)
	letterTemplateRepo := letterTemplateRepository.NewLetterTemplateRepositoryImpl(db)
	letterTemplateSvc := letterTemplateService.NewLetterTemplateServiceImpl(letterTemplateRepo, validate)

	// Seed sla_max_days default if not present, for the default org seeded by
	// db/migrations/002_add_organization_id.sql.
	const defaultOrgID = "00000000-0000-0000-0000-000000000001"
	if existing, _ := appSettingsRepository.GetByKey("sla_max_days", defaultOrgID); existing == nil {
		appSettingsRepository.Update([]model.AppSettings{{Key: "sla_max_days", Value: "7"}}, defaultOrgID)
	}

	// Cron: auto-approve SLA-exceeded documents every day at 06:00
	c := cron.New()
	c.AddFunc("0 6 * * *", func() { slaSvc.RunAutoApprove() })
	c.Start()

	// Controllers
	userController := controller.NewUserController(userService, userLogService)
	documentController := controller.NewDocumentController(documentService, documentNumbersService, userLogService)
	documentHistoryController := controller.NewDocumentHistoryController(documentHistoryService)
	documentSequenceController := controller.NewDocumentSequenceController(documentSequenceService)
	documentAttachmentController := controller.NewDocumentAttachmentController(documentAttachmentService, userLogService)
	positionController := controller.NewPositionController(positionService, userLogService)
	userLogController := controller.NewUserLogController(userLogService)
	appSettingsController := controller.NewAppSettingsController(appSettingsService, userLogService)
	recipientController := controller.NewRecipientController(recipientService)
	bookmarkController := controller.NewBookmarkController(bookmarkService)
	numberingGroupController := controller.NewNumberingGroupController(numberingGroupService, userLogService)
	numberingFormatController := controller.NewNumberingFormatController(numberingFormatService, userLogService)
	documentNumbersController := controller.NewDocumentNumbersController(documentNumbersService, userLogService)
	signatureController := controller.NewSignatureController(signatureService)
	delegatorController := controller.NewDelegatorController(delegatorSvc, userLogService)
	verificationController := controller.NewVerificationController(documentService)
	letterTemplateController := controller.NewLetterTemplateController(letterTemplateSvc, userLogService)

	// S3 presigner (AUDIT SEC-01): the backend now holds the S3 credentials and
	// hands the browser short-lived presigned URLs. New(...) returns nil when
	// creds are absent; the controller reports 503 in that case.
	s3Presigner := s3presign.New(envConf.S3Region, envConf.S3Bucket, envConf.S3AccessKeyID, envConf.S3SecretAccessKey)
	if !s3Presigner.Enabled() {
		log.Println("warning: S3 presigner not configured (S3_* env missing) — file uploads will return 503")
	}
	uploadController := controller.NewUploadController(s3Presigner)

	// Initialize Router
	routes := router.NewRouter(
		db,
		jwksClient,
		userController,
		documentController,
		documentHistoryController,
		documentAttachmentController,
		documentSequenceController,
		positionController,
		userLogController,
		appSettingsController,
		recipientController,
		bookmarkController,
		numberingGroupController,
		numberingFormatController,
		documentNumbersController,
		signatureController,
		delegatorController,
		verificationController,
		letterTemplateController,
		uploadController,
		envConf.CORSAllowedOrigins,
	)

	// Intialize Server
	server := &http.Server{
		Addr:           ":" + envConf.ServerPort,
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// AUDIT SEC-10: run the server in a goroutine and block on OS signals for a
	// graceful shutdown. The previous `server.ListenAndServe().Error()` would
	// panic if ListenAndServe ever returned nil, and left no clean shutdown path.
	go func() {
		log.Printf("Message: Server Successfully Running on :%s", envConf.ServerPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Message: Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("Message: Server stopped cleanly")
}

/*
 Redis should store the User ID, with Key Token.ID or Token
 So when log out, we can extracting the Token.ID or Token as a Key, and get the User ID value for Query processing
 Task:
 - Change Key to Token ID / Token Instead of UserID, in Login and RefreshToken
 -
 Open code for Extracting Token ID value from Redis to get User ID in Logout Function
*/

/*
ChangeLog:
- All Redish process in AUTH flow was removed, because we no longer need it to store an identifier
- Add refresh token expired handler, so when refresh token is expired, logout the user from system
- id in payload access token & refresh token is user id
*/
