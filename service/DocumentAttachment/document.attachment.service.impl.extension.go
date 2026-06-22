package documentAttachment

import (
	response "Microservice/data/response/DocumentAttachment"
	"Microservice/model"
)

func (t DocumentAttachmentServiceImpl) mapDocumentAttachmentToDocumentAttachmentResponse(documentHistories []model.DocumentAttachment) []response.DocumentAttachmentResponse {
	responseDocuments := make([]response.DocumentAttachmentResponse, len(documentHistories))
	for i, documentAttachment := range documentHistories {
		responseDocuments[i] = t.convertDocumentAttachmentToDocumentAttachmentResponse(documentAttachment)
	}
	return responseDocuments
}

func (t DocumentAttachmentServiceImpl) convertDocumentAttachmentToDocumentAttachmentResponse(documentAttachment model.DocumentAttachment) response.DocumentAttachmentResponse {
	// Perform necessary conversion logic here, potentially selecting specific fields
	responseDocument := response.DocumentAttachmentResponse{
		Id:           documentAttachment.ID,
		DocumentID:   &documentAttachment.DocumentID,
		OriginalName: documentAttachment.OriginalName,
		FileName:     documentAttachment.FileName,
		Path:         documentAttachment.Path,
		Size:         documentAttachment.Size,
		Type:         documentAttachment.Type,
	}

	return responseDocument
}
