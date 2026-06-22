package controller

import (
	request "Microservice/data/request/Document"
	documentNumberRequest "Microservice/data/request/DocumentNumbers"
	"Microservice/helper"
	"Microservice/helper/enums"
	"Microservice/model"
	service "Microservice/service/Document"
	documentNumberService "Microservice/service/DocumentNumbers"
	userLogService "Microservice/service/UserLog"
	"Microservice/utils"

	"github.com/gin-gonic/gin"
)

type DocumentController struct {
	documentService       service.DocumentService
	documentNumberService documentNumberService.DocumentNumbersService
	userLogService        userLogService.UserLogService
}

func NewDocumentController(service service.DocumentService, documentNumberService documentNumberService.DocumentNumbersService, userLogService userLogService.UserLogService) *DocumentController {
	return &DocumentController{documentService: service, documentNumberService: documentNumberService, userLogService: userLogService}
}

func (controller *DocumentController) Get(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	stringId := ctx.Param("id")

	documentResponse, errDocumentResponse := controller.documentService.GetDocument(stringId, orgID)
	if errDocumentResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentResponse)
		return
	}
	utils.SuccessResponse(ctx, documentResponse)
}

func (controller *DocumentController) GetDetailPreview(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	stringId := ctx.Param("id")

	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	documentResponse, errDocumentResponse := controller.documentService.GetDetailDocument(stringId, *id, orgID)
	if errDocumentResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentResponse)
		return
	}
	utils.SuccessResponse(ctx, documentResponse)
}

func (controller *DocumentController) GetDetailForEdit(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	stringId := ctx.Param("id")

	documentResponse, errDocumentResponse := controller.documentService.GetDetailForEdit(stringId, orgID)
	if errDocumentResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentResponse)
		return
	}
	utils.SuccessResponse(ctx, documentResponse)
}

func (controller *DocumentController) GetAll(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	documentResponse, errDocumentResponse := controller.documentService.GetAllDocument(orgID)
	if errDocumentResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentResponse)
		return
	}
	utils.SuccessResponse(ctx, documentResponse)
	utils.SetCache(ctx, "All Documents", documentResponse)
}

func (controller *DocumentController) GetAllReferences(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	querySubject := ctx.Param("q")
	documentResponse, errDocumentResponse := controller.documentService.GetAllReferences(querySubject, orgID)
	if errDocumentResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentResponse)
		return
	}
	utils.SuccessResponse(ctx, documentResponse)
	utils.SetCache(ctx, "All Reference Documents", documentResponse)
}

func (controller *DocumentController) GetAllAuthorization(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	documentResponse, errDocumentResponse := controller.documentService.GetAllAuthorization(*id, orgID)
	if errDocumentResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentResponse)
		return
	}
	utils.SuccessResponse(ctx, documentResponse)
	utils.SetCache(ctx, "All Documents", documentResponse)
}

func (controller *DocumentController) GetAllInProgress(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	documentResponse, errDocumentResponse := controller.documentService.GetAllInProgress(*id, orgID)
	if errDocumentResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentResponse)
		return
	}
	utils.SuccessResponse(ctx, documentResponse)
	utils.SetCache(ctx, "All Documents", documentResponse)
}

func (controller *DocumentController) GetAllRejected(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	documentResponse, errDocumentResponse := controller.documentService.GetRejectedByAuthorID(*id, orgID)
	if errDocumentResponse != nil {
		utils.ErrorResponse(ctx, *errDocumentResponse)
		return
	}
	utils.SuccessResponse(ctx, documentResponse)
	utils.SetCache(ctx, "All Documents", documentResponse)
}

func (controller *DocumentController) Create(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.CreateDocumentRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	userId, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	newDocument, err := controller.documentService.Create(payload, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	switch payload.PublicationNumberType {
	case 1:
		if payload.PublicationValue == nil {
			utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "PublicationValue is required for auto-generated numbers."})
			return
		}
		docNumReq := documentNumberRequest.DocumentNumbersRequest{NumberingFormatID: *payload.PublicationValue}
		if err := controller.documentNumberService.Create(docNumReq, *userId, newDocument, enums.Saved, orgID); err != nil {
			utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Document Number Request Structure."})
			return
		}
	case 2:
		if payload.PublicationValue == nil {
			utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "PublicationValue is required for booking numbers."})
			return
		}
		if err := controller.documentNumberService.Update(*payload.PublicationValue, newDocument, enums.Saved, orgID); err != nil {
			utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Document Number Request Structure."})
			return
		}
	}

	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Create),
			Module: string(enums.Document),
			Log:    helper.ToJSON(payload),
		},
		orgID,
	)

	utils.SuccessResponse(ctx, nil)
}

func (controller *DocumentController) Update(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.UpdateDocumentRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	userId, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	// Check err immediately before using the returned document
	document, err := controller.documentService.Update(payload, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	isDocumentNumberStored, errDocID := controller.documentNumberService.GetByDocumentID(document.ID, orgID)
	if errDocID != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	if isDocumentNumberStored == nil {
		switch payload.PublicationNumberType {
		case 1:
			if payload.PublicationValue == nil {
				utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "PublicationValue is required for auto-generated numbers."})
				return
			}
			docNumReq := documentNumberRequest.DocumentNumbersRequest{NumberingFormatID: *payload.PublicationValue}
			if err := controller.documentNumberService.Create(docNumReq, *userId, document, enums.Saved, orgID); err != nil {
				utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Document Number Request Structure."})
				return
			}
		case 2:
			if payload.PublicationValue == nil {
				utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "PublicationValue is required for booking numbers."})
				return
			}
			if err := controller.documentNumberService.Update(*payload.PublicationValue, document, enums.Saved, orgID); err != nil {
				utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Document Number Request Structure."})
				return
			}
		}
	}

	controller.userLogService.CreateLog(
		model.UserLog{
			UserID: *helper.GetUserUUID(ctx),
			Action: string(enums.Update),
			Module: string(enums.Document),
			Log:    helper.ToJSON(payload),
		},
		orgID,
	)

	utils.SuccessResponse(ctx, nil)
}

func (controller *DocumentController) Authorize(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload request.Authorize
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	userId, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	if err := controller.documentService.Authorize(payload, *userId, orgID); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	action := enums.Approve
	if payload.State == 2 {
		action = enums.Reject
	} else if payload.State == 3 {
		action = enums.Cancel
	}
	controller.userLogService.CreateLog(model.UserLog{
		UserID: *helper.GetUserUUID(ctx),
		Action: string(action),
		Module: string(enums.Document),
	}, orgID)

	utils.SuccessResponse(ctx, nil)
}

func (controller *DocumentController) GetComplete(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	documentResponses, err := controller.documentService.GetCompleteByAuthorID(*id, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, documentResponses)
}

func (controller *DocumentController) GetDraft(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	documentResponses, err := controller.documentService.GetDraftByAuthorID(*id, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, documentResponses)
}

func (controller *DocumentController) GetAllInbox(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	documentResponses, err := controller.documentService.GetAllInbox(*id, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, documentResponses)
}

func (controller *DocumentController) GetDashboardSummary(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	period := ctx.DefaultQuery("period", "all")

	validPeriods := map[string]bool{"all": true, "today": true, "week": true, "month": true}
	if !validPeriods[period] {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid period. Use: all, today, week, month"})
		return
	}

	summaryResponse, err := controller.documentService.GetDashboardSummary(*id, period, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	utils.SuccessResponse(ctx, summaryResponse)
}

func (controller *DocumentController) GetDeadlines(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	deadlineResponse, err := controller.documentService.GetDeadlines(*id, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	utils.SuccessResponse(ctx, deadlineResponse)
}

func (controller *DocumentController) GetRecentActivities(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	activitiesResponse, err := controller.documentService.GetRecentActivities(*id, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	utils.SuccessResponse(ctx, activitiesResponse)
}

func (controller *DocumentController) GetRecentDocuments(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	id, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	// 0 = semua, 1 = internal, 2 = external
	docTypeStr := ctx.DefaultQuery("type", "0")
	docType := 0
	switch docTypeStr {
	case "1":
		docType = 1
	case "2":
		docType = 2
	}

	recentResponse, err := controller.documentService.GetRecentDocuments(*id, docType, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	utils.SuccessResponse(ctx, recentResponse)
}

func (controller *DocumentController) Recall(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	documentId := ctx.Param("id")
	if documentId == "" {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Document ID is required"})
		return
	}

	userId, errParse := helper.GetUserId(ctx)
	if errParse != nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Invalid Request Structure."})
		return
	}

	if err := controller.documentService.Recall(documentId, *userId, orgID); err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	utils.SuccessResponse(ctx, nil)
}

func (controller *DocumentController) Search(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	keyword := ctx.Query("q")
	if keyword == "" {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Keyword tidak boleh kosong."})
		return
	}

	if len(keyword) > 100 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Keyword terlalu panjang, maksimal 100 karakter."})
		return
	}

	result, err := controller.documentService.Search(keyword, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	utils.SuccessResponse(ctx, result)
}
