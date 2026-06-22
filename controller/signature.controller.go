package controller

import (
	signatureRequest "Microservice/data/request/Signature"
	"Microservice/helper"
	signatureService "Microservice/service/Signature"
	"Microservice/utils"

	"github.com/gin-gonic/gin"
)

type SignatureController struct {
	signatureService signatureService.SignatureService
}

func NewSignatureController(service signatureService.SignatureService) *SignatureController {
	return &SignatureController{signatureService: service}
}

func (controller *SignatureController) Create(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	var payload signatureRequest.CreateSignatureRequest

	errBindJSON := ctx.ShouldBindJSON(&payload)
	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	err := controller.signatureService.Create(payload, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	utils.SuccessResponse(ctx, nil)
}

func (controller *SignatureController) Update(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	userId := ctx.Param("userId")
	var payload signatureRequest.UpdateSignatureRequest

	errBindJSON := ctx.ShouldBindJSON(&payload)
	if errBindJSON != nil {
		msg := "Bad Request"
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: msg})
		return
	}

	if errs := helper.ValidateStruct(payload); len(errs) > 0 {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 400, Message: "Bad Request"})
		return
	}

	err := controller.signatureService.Update(userId, payload, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	utils.SuccessResponse(ctx, nil)
}

func (controller *SignatureController) Delete(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	userId := ctx.Param("userId")

	err := controller.signatureService.Delete(userId, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	utils.SuccessResponse(ctx, nil)
}

func (controller *SignatureController) GetAll(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}

	signatures, err := controller.signatureService.GetAll(orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	utils.SuccessResponse(ctx, signatures)
}

func (controller *SignatureController) GetByUserId(ctx *gin.Context) {
	orgID, ok := helper.RequireOrgID(ctx)
	if !ok {
		return
	}
	userId := ctx.Param("userId")

	signature, err := controller.signatureService.GetByUserId(userId, orgID)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}

	if signature == nil {
		utils.ErrorResponse(ctx, helper.ErrorModel{Code: 404, Message: "Signature not found"})
		return
	}

	utils.SuccessResponse(ctx, signature)
}
