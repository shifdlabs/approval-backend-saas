package controller

import (
	service "Microservice/service/Document"
	"Microservice/utils"

	"github.com/gin-gonic/gin"
)

type VerificationController struct {
	documentService service.DocumentService
}

func NewVerificationController(service service.DocumentService) *VerificationController {
	return &VerificationController{documentService: service}
}

func (controller *VerificationController) GetVerification(ctx *gin.Context) {
	documentId := ctx.Param("id")

	result, err := controller.documentService.GetVerification(documentId)
	if err != nil {
		utils.ErrorResponse(ctx, *err)
		return
	}
	utils.SuccessResponse(ctx, result)
}
