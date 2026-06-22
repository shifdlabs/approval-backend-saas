package document

import (
	request "Microservice/data/request/Document"
	response "Microservice/data/response/Document"
	"Microservice/helper"
	"Microservice/model"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

func (t DocumentServiceImpl) mapDocumentsToDocumentResponse(documents []model.Document) []response.DocumentResponse {
	responseDocuments := make([]response.DocumentResponse, len(documents))
	for i, document := range documents {
		responseDocuments[i] = t.convertDocumentToDocumentResponse(document)
	}
	return responseDocuments
}

func (t DocumentServiceImpl) convertDocumentToDocumentResponse(document model.Document) response.DocumentResponse {

	var currentApproverTitle *string
	var lastRejector *response.RejectorResponse

	if len(document.DocumentSequence) > 0 {
		currentApprover, _ := t.DocumentSequenceRepository.GetCurrentApprover(document.ID.String())

		if currentApprover != nil {
			user, _ := t.UserRepository.Get(currentApprover.UserID.String(), true, document.OrganizationID.String())

			if user != nil {
				positionName := ""
				if user.Position != nil {
					positionName = user.Position.Name
				}
				title := fmt.Sprintf("%s %s - %s",
					user.FirstName,
					user.LastName,
					positionName,
				)

				currentApproverTitle = &title
			}
		}
	}

	lastRejectorResponse, _ := t.DocumentHistoryRepository.GetLastRejection(document.ID.String())
	if lastRejectorResponse != nil {
		rejectorData, _ := t.UserRepository.Get(string(lastRejectorResponse.UserID.String()), true, document.OrganizationID.String())

		if rejectorData != nil {
			rejectorName := fmt.Sprintf("%s %s",
				rejectorData.FirstName,
				rejectorData.LastName,
			)

			lastRejector = &response.RejectorResponse{
				Name:   &rejectorName,
				Reason: &lastRejectorResponse.Description,
			}
		}
	}

	// Perform necessary conversion logic here, potentially selecting specific fields
	responseDocument := response.DocumentResponse{
		Id:                  &document.ID,
		Subject:             document.Subject,
		Body:                document.Body,
		Type:                document.Type,
		Step:                document.Step,
		Status:              document.Status,
		Priority:            document.Priority,
		Author:              document.Author,
		DocumentSequence:    document.DocumentSequence,
		DocumentHistory:     document.DocumentHistory,
		DocumentAttachment:  document.DocumentAttachment,
		CreatedAt:           *document.CreatedAt,
		UpdatedAt:           *document.UpdatedAt,
		CurrentApprovalName: currentApproverTitle,
		LastRejector:        lastRejector,
	}

	return responseDocument
}

func (t DocumentServiceImpl) convertToDocumentDetailResponse(document model.Document, userId string) response.DocumentDetailResponse {
	orgID := document.OrganizationID.String()
	inProgressOverview, _ := t.GetInProgressOverviewByDocId(document.ID.String(), orgID)

	documentSequence, _ := t.DocumentSequenceRepository.GetSequencesByDocumentId(document.ID.String())
	currentApprover := t.getCurrentApprover(documentSequence, document)
	documentHistories := t.getDocumentHistory(document)
	documentAttachment := t.getDocumentAttachment(document)
	recipients := t.getInternalRecipients(document.ID.String(), orgID)
	referencesResult, _ := t.DocumentReferenceRepository.GetAll(document.ID)
	documentReferences := make([]response.DocumentReference, len(referencesResult))

	if len(referencesResult) > 0 {
		for i, reference := range referencesResult {
			helper.PrintValue(reference.DocumentID, "Reference ID")
			document, err := t.DocumentRepository.Get(reference.DocumentID.String(), orgID)
			if err != nil || document == nil {
				// Skip this reference, move on to the next
				continue
			}

			documentReferences[i] = response.DocumentReference{
				Id:      document.ID.String(),
				Subject: document.Subject,
			}
		}
	}

	var publicationValue string
	switch document.PublicationNumberType {
	case 1, 2:
		response, _ := t.DocumentNumbersRepository.GetByDocumentID(document.ID, orgID)
		if response != nil {
			publicationValue = response.Value
		}
	case 3:
		publicationValue = *document.CustomPublicationNumber

	case 4:
		publicationValue = ""
	}

	var author model.User
	if document.Author != nil {
		author = *document.Author
	}
	isAllowToUpdate := document.Author != nil && document.Author.ID.String() == userId && document.Status == 99
	canRecall := document.Author != nil && document.Author.ID.String() == userId && document.Status == 1 && len(document.DocumentHistory) == 0

	response := response.DocumentDetailResponse{
		Id:                 &document.ID,
		PublicationValue:   publicationValue,
		ExternalRecipient:  document.ExternalRecipient,
		Subject:            document.Subject,
		Body:               document.Body,
		Type:               document.Type,
		Step:               document.Step,
		Status:             document.Status,
		Priority:           document.Priority,
		Author:             author,
		DocumentSequence:   *inProgressOverview,
		DocumentHistory:    &documentHistories,
		DocumentAttachment: &documentAttachment,
		DocumentReferences: &documentReferences,
		InternalRecipients: &recipients,
		CreatedAt:          *document.CreatedAt,
		UpdatedAt:          *document.UpdatedAt,
		IsApprover:         t.isApproverOrDelegate(currentApprover, document.Status, userId, orgID),
		IsAllowToUpdate:    isAllowToUpdate,
		CanRecall:          canRecall,
		DueDate:            document.DueDate,
	}

	return response
}

func (t DocumentServiceImpl) getInternalRecipients(documentId string, orgID string) []response.InternalRecipient {
	recipientsResponse, _ := t.RecipientRepository.GetRecipientsByDocId(documentId)
	recipients := make([]response.InternalRecipient, 0, len(recipientsResponse))
	for _, recipient := range recipientsResponse {
		user, _ := t.UserRepository.Get(recipient.UserID.String(), true, orgID)
		if user == nil {
			continue
		}
		positionName := ""
		if user.Position != nil {
			positionName = user.Position.Name
		}
		recipients = append(recipients, response.InternalRecipient{
			Name:  user.FirstName + " " + user.LastName,
			Title: positionName,
		})
	}

	return recipients
}

func (t DocumentServiceImpl) getDocumentAttachment(document model.Document) []response.DocumentAttachment {
	documentAttachment := make([]response.DocumentAttachment, len(document.DocumentAttachment))
	for i, attachment := range document.DocumentAttachment {
		documentAttachment[i] = response.DocumentAttachment{
			Id:           attachment.ID.String(),
			OriginalName: attachment.OriginalName,
			FileName:     attachment.FileName,
			Path:         attachment.Path,
			Size:         attachment.Size,
			Type:         attachment.Type,
		}
	}

	return documentAttachment
}

func (t DocumentServiceImpl) getDocumentHistory(document model.Document) []response.DocumentHistory {
	orgID := document.OrganizationID.String()
	documentHistories := make([]response.DocumentHistory, len(document.DocumentHistory))
	for i, history := range document.DocumentHistory {
		user, _ := t.UserRepository.Get(history.UserID.String(), true, orgID)
		name := ""
		positionName := ""
		if user != nil {
			name = user.FirstName + " " + user.LastName
			if user.Position != nil {
				positionName = user.Position.Name
			}
		}
		documentHistories[i] = response.DocumentHistory{
			Name:       name,
			Title:      positionName,
			IsApproved: history.IsApproved,
			Reason:     history.Description,
			UpdatedAt:  history.CreatedAt.String(),
		}
	}

	return documentHistories
}

func (t DocumentServiceImpl) getApprovers(document model.Document) []string {
	orgID := document.OrganizationID.String()
	approverIds := make([]string, len(document.DocumentSequence))
	for i, history := range document.DocumentSequence {
		user, _ := t.UserRepository.Get(history.UserID.String(), true, orgID)

		approverIds[i] = user.ID.String()
	}

	return approverIds
}

// isApproverOrDelegate returns true when userId is either the direct sequence approver
// or is in the owner-chain that delegates authority to userId on today's date.
func (t DocumentServiceImpl) isApproverOrDelegate(currentApprover model.DocumentSequence, documentStatus int, userId string, orgID string) bool {
	if documentStatus != 1 {
		fmt.Printf("[isApproverOrDelegate] status=%d != 1, returning false\n", documentStatus)
		return false
	}
	seqUserID := currentApprover.UserID.String()
	fmt.Printf("[isApproverOrDelegate] seqUserID=%s userId=%s\n", seqUserID, userId)
	if seqUserID == userId {
		return true
	}
	ownerChain, err := t.getOwnerChainForDelegate(userId, time.Now(), orgID)
	fmt.Printf("[isApproverOrDelegate] ownerChain=%v err=%v\n", ownerChain, err)
	if err != nil {
		return false
	}
	for _, ownerID := range ownerChain {
		if ownerID == seqUserID {
			fmt.Printf("[isApproverOrDelegate] delegate match found, returning true\n")
			return true
		}
	}
	return false
}

func (t DocumentServiceImpl) getCurrentApprover(sequences []model.DocumentSequence, document model.Document) model.DocumentSequence {
	currentApprover := model.DocumentSequence{}
	if len(sequences) > 0 {
		currentApprover = sequences[document.Step-1]
	}

	return currentApprover
}

func (t DocumentServiceImpl) convertRequestToCreateModel(documentRequest request.CreateDocumentRequest, user *model.User) (*model.Document, *helper.ErrorModel) {
	var customPublicationCode *string = nil
	if documentRequest.PublicationNumberType == 3 {
		customPublicationCode = documentRequest.PublicationValue
	}

	var templateUUID *uuid.UUID
	if documentRequest.TemplateID != nil && *documentRequest.TemplateID != "" {
		if parsed, err := uuid.FromString(*documentRequest.TemplateID); err == nil {
			templateUUID = &parsed
		}
	}

	// Store to DB
	document := model.Document{
		Author:                  user,
		PublicationNumberType:   documentRequest.PublicationNumberType,
		CustomPublicationNumber: customPublicationCode,
		Type:                    documentRequest.Type,
		Priority:                documentRequest.Priority,
		Subject:                 documentRequest.Subject,
		Body:                    documentRequest.Body,
		ExternalRecipient:       documentRequest.ExternalRecipient,
		Step:                    documentRequest.Step,
		LetterHead:              documentRequest.LetterHead,
		Status:                  documentRequest.Status,
		TemplateID:              templateUUID,
	}

	return &document, nil
}

func (t DocumentServiceImpl) convertDocumentToEditDocumentResponse(document model.Document) response.EditDocumentResponse {

	internalRecipient, _ := t.RecipientRepository.GetRecipientsByDocId(document.ID.String())
	carbonCopy, _ := t.CarbonCopyRepository.GetCarbonCopysByDocId(document.ID.String())
	documentAttachment := t.getDocumentAttachment(document)
	approvers := t.getApprovers(document)
	recipientIds := make([]string, 0, len(internalRecipient))
	for _, r := range internalRecipient {
		recipientIds = append(recipientIds, string(r.UserID.String()))
	}

	carbonCopiesIds := make([]string, 0, len(carbonCopy))
	for _, r := range carbonCopy {
		carbonCopiesIds = append(carbonCopiesIds, string(r.UserID.String()))
	}

	referencesResult, _ := t.DocumentReferenceRepository.GetAll(document.ID)
	documentReferences := make([]response.DocumentReference, len(referencesResult))

	if len(referencesResult) > 0 {
		orgID := document.OrganizationID.String()
		for i, reference := range referencesResult {
			document, err := t.DocumentRepository.Get(reference.DocumentID.String(), orgID)
			if err != nil || document == nil {
				// Skip this reference, move on to the next
				continue
			}

			documentReferences[i] = response.DocumentReference{
				Id:      document.ID.String(),
				Subject: document.Subject,
			}
		}
	}

	var publicationNumber string

	switch document.PublicationNumberType {
	case 1, 2:
		documentNumber, _ := t.DocumentNumbersRepository.GetByDocumentID(document.ID, document.OrganizationID.String())
		publicationNumber = documentNumber.Value
	case 3:
		publicationNumber = *document.CustomPublicationNumber
	default:
		publicationNumber = ""
	}

	response := response.EditDocumentResponse{
		Id:                    &document.ID,
		PublicationNumberType: document.PublicationNumberType,
		PublicationValue:      &publicationNumber,
		DocumentReferences:    &documentReferences,
		Subject:               document.Subject,
		Body:                  document.Body,
		Type:                  document.Type,
		Step:                  document.Step,
		Status:                document.Status,
		Priority:              document.Priority,
		Author:                *document.Author,
		DocumentAttachment:    &documentAttachment,
		ExternalRecipient:     &document.ExternalRecipient,
		InternalRecipients:    &recipientIds,
		CarbonCopy:            &carbonCopiesIds,
		Approvers:             &approvers,
	}

	return response
}
