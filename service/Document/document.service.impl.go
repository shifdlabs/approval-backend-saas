package document

import (
	request "Microservice/data/request/Document"
	response "Microservice/data/response/Document"
	"Microservice/helper"
	"Microservice/model"
	appSettingsRepository "Microservice/repository/AppSettings"
	carbonCopyReposiory "Microservice/repository/CarbonCopy"
	delegatorRepository "Microservice/repository/Delegator"
	repository "Microservice/repository/Document"
	documentAttachmentRepository "Microservice/repository/DocumentAttachment"
	documentHistoryReposiory "Microservice/repository/DocumentHistory"
	documentNumbersRepository "Microservice/repository/DocumentNumbers"
	documentReferenceRepository "Microservice/repository/DocumentReference"
	documentSequenceReposiory "Microservice/repository/DocumentSequence"
	recipientReposiory "Microservice/repository/Recipient"
	signatureRepository "Microservice/repository/Signature"
	userRepository "Microservice/repository/User"
	userLogRepository "Microservice/repository/UserLog"
	emailService "Microservice/service/Email"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type DocumentServiceImpl struct {
	DocumentRepository           repository.DocumentRepository
	UserRepository               userRepository.UserRepository
	DocumentSequenceRepository   documentSequenceReposiory.DocumentSequenceRepository
	DocumentAttachmentRepository documentAttachmentRepository.DocumentAttachmentRepository
	DocumentHistoryRepository    documentHistoryReposiory.DocumentHistoryRepository
	RecipientRepository          recipientReposiory.RecipientRepository
	CarbonCopyRepository         carbonCopyReposiory.CarbonCopyRepository
	UserLogRepository            userLogRepository.UserLogRepository
	DocumentNumbersRepository    documentNumbersRepository.DocumentNumbersRepository
	DocumentReferenceRepository  documentReferenceRepository.DocumentReferenceRepository
	SignatureRepository          signatureRepository.SignatureRepository
	DelegatorRepository          delegatorRepository.DelegatorRepository
	AppSettingsRepository        appSettingsRepository.AppSettingsRepository
	EmailService                 emailService.EmailService
	FrontendURL                  string
	Validate                     *validator.Validate
	Db                           *gorm.DB
}

func NewDocumentServiceImpl(
	documentRepository repository.DocumentRepository,
	userRepository userRepository.UserRepository,
	documentSequenceRepository documentSequenceReposiory.DocumentSequenceRepository,
	documentAttachmentRepository documentAttachmentRepository.DocumentAttachmentRepository,
	documentHistoryRepository documentHistoryReposiory.DocumentHistoryRepository,
	recipientRepository recipientReposiory.RecipientRepository,
	carbonCopyRepository carbonCopyReposiory.CarbonCopyRepository,
	userLogRepository userLogRepository.UserLogRepository,
	documentNumbersRepository documentNumbersRepository.DocumentNumbersRepository,
	documentReferenceRepository documentReferenceRepository.DocumentReferenceRepository,
	signatureRepository signatureRepository.SignatureRepository,
	delegatorRepo delegatorRepository.DelegatorRepository,
	appSettingsRepo appSettingsRepository.AppSettingsRepository,
	emailSvc emailService.EmailService,
	frontendURL string,
	Db *gorm.DB,
	validate *validator.Validate) DocumentService {
	return &DocumentServiceImpl{
		DocumentRepository:           documentRepository,
		UserRepository:               userRepository,
		DocumentSequenceRepository:   documentSequenceRepository,
		DocumentAttachmentRepository: documentAttachmentRepository,
		DocumentHistoryRepository:    documentHistoryRepository,
		RecipientRepository:          recipientRepository,
		CarbonCopyRepository:         carbonCopyRepository,
		UserLogRepository:            userLogRepository,
		DocumentNumbersRepository:    documentNumbersRepository,
		DocumentReferenceRepository:  documentReferenceRepository,
		SignatureRepository:          signatureRepository,
		DelegatorRepository:          delegatorRepo,
		AppSettingsRepository:        appSettingsRepo,
		EmailService:                 emailSvc,
		FrontendURL:                  frontendURL,
		Db:                           Db,
		Validate:                     validate,
	}
}

func (t DocumentServiceImpl) Create(request request.CreateDocumentRequest, orgID string) (*model.Document, *helper.ErrorModel) {
	errStructure := t.Validate.Struct(request)
	if errStructure != nil {
		msg := "Structure Error"
		return nil, helper.ErrorCatcher(errStructure, 500, &msg)
	}

	orgUUID, errParseOrg := uuid.FromString(orgID)
	if errParseOrg != nil {
		msg := "Invalid Organization ID"
		return nil, helper.ErrorCatcher(errParseOrg, 500, &msg)
	}

	// Get User Data
	user, errUser := t.UserRepository.Get(request.AuthorID, true, orgID)
	if errUser != nil {
		return nil, nil
	}

	// Start Transaction
	trx := t.Db.Begin()

	newDocument, errConvert := t.convertRequestToCreateModel(request, user)
	if errConvert != nil {
		return nil, errConvert
	}
	newDocument.OrganizationID = &orgUUID

	t.DocumentRepository.Create(*trx, newDocument)

	if newDocument.ID == uuid.Nil {
		trx.Rollback()
		msg := "Document ID is nil after creation"
		return nil, helper.ErrorCatcher(errStructure, 500, &msg)
	}

	// Store Internal Recipients
	if request.Recipients != nil {
		var recipients []model.Recipient
		for _, userId := range request.Recipients {
			userUuid, errorParse := uuid.FromString(userId)
			if errorParse != nil {
				trx.Rollback()
				msg := "Structure Error"
				return nil, helper.ErrorCatcher(errStructure, 500, &msg)
			}

			recipients = append(recipients, model.Recipient{
				Document: newDocument,
				UserID:   userUuid,
			})
		}

		t.RecipientRepository.Create(
			*trx,
			recipients,
		)
	}

	// Store Carbon Copies
	if request.CarbonCopies != nil {
		var carbonCopies []model.CarbonCopy
		for _, userId := range request.CarbonCopies {
			userUuid, errorParse := uuid.FromString(userId)
			if errorParse != nil {
				trx.Rollback()
				msg := "Structure Error"
				return nil, helper.ErrorCatcher(errStructure, 500, &msg)
			}

			carbonCopies = append(carbonCopies, model.CarbonCopy{
				Document: newDocument,
				UserID:   userUuid,
			})
		}

		t.CarbonCopyRepository.Create(
			*trx,
			carbonCopies,
		)
	}

	// Store Document Approvers
	for index, value := range request.Sequences {
		userId, _ := uuid.FromString(value.UserID)
		t.DocumentSequenceRepository.Create(
			trx,
			model.DocumentSequence{
				DocumentID: &newDocument.ID,
				UserID:     userId,
				Step:       (index + 1),
				Signature:  value.Signature,
			},
		)
	}

	// Store Document Sequences
	for _, value := range request.Attachments {
		t.DocumentAttachmentRepository.Create(
			trx,
			model.DocumentAttachment{
				Document:     newDocument,
				OriginalName: value.OriginalName,
				FileName:     value.FileName,
				Path:         value.Path,
				Size:         value.Size,
				Type:         value.Type,
			},
		)
	}

	if request.References != nil {
		for _, referenceID := range request.References {
			referenceUUID, _ := uuid.FromString(referenceID)
			t.DocumentReferenceRepository.Create(trx, model.DocumentReference{
				ReferenceID: referenceUUID,
				DocumentID:  newDocument.ID,
			})
		}
	}

	// End Transaction
	trx.Commit()

	// Notify first approver if document is submitted (not a draft)
	if newDocument.Status == 1 && len(request.Sequences) > 0 {
		docID := newDocument.ID.String()
		docSubject := newDocument.Subject
		firstApproverID := request.Sequences[0].UserID
		authorID := request.AuthorID
		frontendURL := t.FrontendURL
		go func() {
			approver, err := t.UserRepository.Get(firstApproverID, true, orgID)
			if err != nil || approver == nil {
				return
			}
			author, err := t.UserRepository.Get(authorID, true, orgID)
			if err != nil || author == nil {
				return
			}
			fromName := author.FirstName + " " + author.LastName
			documentURL := fmt.Sprintf("%s/preview/%s", frontendURL, docID)
			t.EmailService.SendApprovalRequest(approver.Email, approver.FirstName+" "+approver.LastName, fromName, docSubject, documentURL)
		}()
	}

	return newDocument, nil
}

func (t DocumentServiceImpl) GetDocument(id string, orgID string) (*response.DocumentResponse, *helper.ErrorModel) {
	document, fetchError := t.DocumentRepository.Get(id, orgID)

	if fetchError != nil {
		return nil, fetchError
	}

	documentResponse := t.convertDocumentToDocumentResponse(*document)
	return &documentResponse, fetchError
}

func (t DocumentServiceImpl) GetDetailDocument(id string, currentUserId string, orgID string) (*response.DocumentDetailResponse, *helper.ErrorModel) {
	helper.PrintValue("Rezz", id)
	document, fetchError := t.DocumentRepository.Get(id, orgID)

	if fetchError != nil {
		return nil, fetchError
	}

	documentResponse := t.convertToDocumentDetailResponse(*document, currentUserId)

	return &documentResponse, fetchError
}

func (t DocumentServiceImpl) GetDetailForEdit(id string, orgID string) (*response.EditDocumentResponse, *helper.ErrorModel) {
	document, fetchError := t.DocumentRepository.Get(id, orgID)

	if fetchError != nil {
		return nil, fetchError
	}

	documentResponse := t.convertDocumentToEditDocumentResponse(*document)

	return &documentResponse, fetchError
}

func (t DocumentServiceImpl) GetAllDocument(orgID string) ([]response.DocumentResponse, *helper.ErrorModel) {
	result, fetchError := t.DocumentRepository.GetAll(orgID)

	if fetchError != nil {
		return nil, fetchError
	} else {
		return t.mapDocumentsToDocumentResponse(result), nil
	}
}

func (t DocumentServiceImpl) GetAllReferences(query string, orgID string) ([]response.DocumentResponse, *helper.ErrorModel) {
	result, fetchError := t.DocumentRepository.GetAllReferences(query, orgID)

	if fetchError != nil {
		return nil, fetchError
	} else {
		return t.mapDocumentsToDocumentResponse(result), nil
	}
}

func (t DocumentServiceImpl) GetAllAuthorization(userId string, orgID string) ([]response.DocumentResponse, *helper.ErrorModel) {
	today := time.Now()

	// Find all ownerIDs whose delegation chain resolves to the current user (BFS backward)
	ownerChain, errChain := t.getOwnerChainForDelegate(userId, today, orgID)
	if errChain != nil {
		return nil, errChain
	}

	userIDs := append([]string{userId}, ownerChain...)

	var allDocuments []model.Document
	seen := map[string]bool{}

	for _, uid := range userIDs {
		docs, fetchError := t.DocumentRepository.GetAllAuthorization(uid, orgID)
		if fetchError != nil {
			return nil, fetchError
		}
		for _, doc := range docs {
			docID := doc.ID.String()
			if !seen[docID] {
				// For delegated docs, verify end-user of the chain is really the current user
				if uid != userId {
					resolved, errResolve := t.resolveDelegate(uid, today, orgID)
					if errResolve != nil || resolved != userId {
						continue
					}
				}
				seen[docID] = true
				allDocuments = append(allDocuments, doc)
			}
		}
	}

	return t.mapDocumentsToDocumentResponse(allDocuments), nil
}

// resolveDelegate follows the delegation chain from ownerID and returns the final delegate ID.
func (t DocumentServiceImpl) resolveDelegate(ownerID string, date time.Time, orgID string) (string, *helper.ErrorModel) {
	const maxDepth = 10
	current := ownerID
	for i := 0; i < maxDepth; i++ {
		delegation, err := t.DelegatorRepository.GetActiveDelegationByOwnerID(current, date, orgID)
		if err != nil {
			return "", err
		}
		if delegation == nil {
			return current, nil
		}
		next := delegation.DelegateID.String()
		if next == ownerID {
			return current, nil // circular — stop
		}
		current = next
	}
	return current, nil
}

// getOwnerChainForDelegate returns all ownerIDs whose delegation chain eventually resolves to delegateID.
func (t DocumentServiceImpl) getOwnerChainForDelegate(delegateID string, date time.Time, orgID string) ([]string, *helper.ErrorModel) {
	visited := map[string]bool{delegateID: true}
	result := []string{}
	queue := []string{delegateID}

	for len(queue) > 0 && len(result) < 100 {
		current := queue[0]
		queue = queue[1:]

		ownerIDs, err := t.DelegatorRepository.GetOwnerIDsByDelegateID(current, date, orgID)
		if err != nil {
			return nil, err
		}
		for _, ownerID := range ownerIDs {
			if !visited[ownerID] {
				visited[ownerID] = true
				result = append(result, ownerID)
				queue = append(queue, ownerID)
			}
		}
	}
	return result, nil
}

func (t DocumentServiceImpl) GetAllInProgress(userId string, orgID string) ([]response.DocumentResponse, *helper.ErrorModel) {
	result, fetchError := t.DocumentRepository.GetAllInProgress(userId, orgID)
	if fetchError != nil {
		return nil, fetchError
	} else {
		return t.mapDocumentsToDocumentResponse(result), nil
	}
}

func (t DocumentServiceImpl) GetDocumentStatistics(userId string, orgID string) (*response.DocumentStatisticResponse, *helper.ErrorModel) {
	result, fetchError := t.DocumentRepository.GetDocumentStatistics(userId, orgID)

	if len(result) == 4 {
		return &response.DocumentStatisticResponse{
			Authorization: result[0],
			InProgress:    result[1],
			Rejected:      result[2],
			Completed:     result[3],
		}, nil
	} else {
		return nil, fetchError
	}
}

func (t DocumentServiceImpl) GetInProgressOverviewByDocId(documentId string, orgID string) (*response.DocumentInProgressResponse, *helper.ErrorModel) {
	// InProgress Overview
	document, fetchDocErr := t.DocumentRepository.Get(documentId, orgID)
	approvers := []response.ApproverForOverview{}

	if fetchDocErr != nil {
		return nil, fetchDocErr
	}

	if document != nil {
		today := time.Now()

		allSequences, _ := t.DocumentSequenceRepository.GetSequencesByDocumentId(document.ID.String())
		sequenceMap := make(map[string]*model.DocumentSequence)
		if allSequences != nil {
			for i, seq := range allSequences {
				sequenceMap[seq.UserID.String()] = &allSequences[i]
			}
		}

		history, fetchHistoryErr := t.DocumentHistoryRepository.GetAllHistoryByDocumentId(document.ID.String())
		if fetchHistoryErr != nil {
			return nil, fetchDocErr
		}

		// addedUserIds tracks which sequence user IDs are already represented in the timeline
		addedUserIds := make(map[string]bool)

		if history != nil {
			for _, approver := range history {
				var displayUser *model.User
				var onBehalfOfName *string

				if approver.OnBehalfOfID != nil {
					// Delegate approval: show the original approver's name, note who acted
					originalUser, err := t.UserRepository.Get(approver.OnBehalfOfID.String(), true, orgID)
					if err != nil {
						return nil, err
					}
					displayUser = originalUser

					delegateUser, err := t.UserRepository.Get(approver.UserID.String(), true, orgID)
					if err != nil {
						return nil, err
					}
					delegateName := delegateUser.FirstName + " " + delegateUser.LastName
					onBehalfOfName = &delegateName

					// Mark the original sequence user as handled (not the delegate)
					addedUserIds[originalUser.ID.String()] = true
					addedUserIds[approver.UserID.String()] = true
				} else {
					user, err := t.UserRepository.Get(approver.UserID.String(), true, orgID)
					if err != nil {
						return nil, err
					}
					displayUser = user
					addedUserIds[user.ID.String()] = true
				}

				var signatureUrl *string
				signature, _ := t.SignatureRepository.GetByUserId(displayUser.ID.String(), orgID)
				if signature != nil {
					signatureUrl = &signature.ImageURL
				}

				hasSigned := false
				if seq, exists := sequenceMap[displayUser.ID.String()]; exists {
					hasSigned = seq.Signature
				}

				positionName := ""
				if displayUser.Position != nil {
					positionName = displayUser.Position.Name
				}

				dateStr := approver.CreatedAt.String()
				approvers = append(approvers, response.ApproverForOverview{
					Name:         displayUser.FirstName + " " + displayUser.LastName,
					Title:        positionName,
					Approved:     &approver.IsApproved,
					Date:         &dateStr,
					Signature:    hasSigned,
					SignatureUrl: signatureUrl,
					OnBehalfOf:   onBehalfOfName,
				})
			}
		}

		if allSequences != nil {
			for _, sequence := range allSequences {
				if addedUserIds[sequence.UserID.String()] {
					continue
				}

				user, err := t.UserRepository.Get(sequence.UserID.String(), true, orgID)
				if err != nil {
					return nil, err
				}

				var signatureUrl *string
				signature, _ := t.SignatureRepository.GetByUserId(user.ID.String(), orgID)
				if signature != nil {
					signatureUrl = &signature.ImageURL
				}

				// Check if this pending approver has an active delegation
				var delegateName *string
				delegation, _ := t.DelegatorRepository.GetActiveDelegationByOwnerID(sequence.UserID.String(), today, orgID)
				if delegation != nil {
					delegateUser, err := t.UserRepository.Get(delegation.DelegateID.String(), true, orgID)
					if err == nil {
						name := delegateUser.FirstName + " " + delegateUser.LastName
						delegateName = &name
					}
				}

				positionName := ""
				if user.Position != nil {
					positionName = user.Position.Name
				}

				approvers = append(approvers, response.ApproverForOverview{
					Name:         user.FirstName + " " + user.LastName,
					Title:        positionName,
					Approved:     nil,
					Date:         nil,
					Signature:    sequence.Signature,
					SignatureUrl: signatureUrl,
					DelegateName: delegateName,
				})
			}
		}
	}

	if document != nil {
		return &response.DocumentInProgressResponse{
			Subject:   document.Subject,
			Approvers: approvers,
		}, nil
	} else {
		return nil, nil
	}
}

func (t DocumentServiceImpl) GetInProgressOverview(userId string, orgID string) (*response.DocumentInProgressResponse, *helper.ErrorModel) {
	// InProgress Overview
	document, fetchDocErr := t.DocumentRepository.GetOneLatestInprogress(userId, orgID)

	approvers := []response.ApproverForOverview{}

	if fetchDocErr != nil {
		return nil, fetchDocErr
	}

	if document != nil {
		allSequences, _ := t.DocumentSequenceRepository.GetSequencesByDocumentId(document.ID.String())
		sequenceMap := make(map[string]*model.DocumentSequence)
		if allSequences != nil {
			for i, seq := range allSequences {
				sequenceMap[seq.UserID.String()] = &allSequences[i]
			}
		}

		history, fetchHistoryErr := t.DocumentHistoryRepository.GetAllHistoryByDocumentId(document.ID.String())
		if fetchHistoryErr != nil {
			return nil, fetchDocErr
		}

		// Track which users are already added from history
		addedUserIds := make(map[string]bool)

		if history != nil {
			for _, approver := range history {
				user, err := t.UserRepository.Get(approver.UserID.String(), true, orgID)
				if err != nil {
					return nil, err
				}

				// Mark this user as added
				addedUserIds[user.ID.String()] = true

				// Fetch signature if exists
				var signatureUrl *string
				signature, _ := t.SignatureRepository.GetByUserId(user.ID.String(), orgID)
				if signature != nil {
					signatureUrl = &signature.ImageURL
				}

				// Get sequence to check if they signed
				hasSigned := false
				if seq, exists := sequenceMap[user.ID.String()]; exists {
					hasSigned = seq.Signature
				}

				positionName := ""
				if user.Position != nil {
					positionName = user.Position.Name
				}

				dateStr := approver.CreatedAt.String()
				approvers = append(approvers, response.ApproverForOverview{
					Name:         user.FirstName + " " + user.LastName,
					Title:        positionName,
					Approved:     &approver.IsApproved,
					Date:         &dateStr,
					Signature:    hasSigned,
					SignatureUrl: signatureUrl,
				})
			}
		}

		if allSequences != nil {
			for _, sequence := range allSequences {

				if addedUserIds[sequence.UserID.String()] {
					continue
				}

				user, err := t.UserRepository.Get(sequence.UserID.String(), true, orgID)
				if err != nil {
					return nil, err
				}

				var signatureUrl *string
				signature, _ := t.SignatureRepository.GetByUserId(user.ID.String(), orgID)
				if signature != nil {
					signatureUrl = &signature.ImageURL
				}

				positionName := ""
				if user.Position != nil {
					positionName = user.Position.Name
				}

				approvers = append(approvers, response.ApproverForOverview{
					Name:         user.FirstName + " " + user.LastName,
					Title:        positionName,
					Approved:     nil,
					Date:         nil,
					Signature:    sequence.Signature,
					SignatureUrl: signatureUrl,
				})
			}
		}
	}

	if document != nil {
		return &response.DocumentInProgressResponse{
			Subject:   document.Subject,
			Approvers: approvers,
		}, nil
	} else {
		return nil, nil
	}
}

func (t DocumentServiceImpl) GetRejectedOverview(userId string, orgID string) (*response.RejectedOverviewResponse, *helper.ErrorModel) {
	document, fetchDocErr := t.DocumentRepository.GetLastestRejected(userId, orgID)

	rejected := response.RejectedOverviewResponse{}

	if fetchDocErr != nil {
		return nil, fetchDocErr
	}

	if document != nil {
		rejectedBy, fetchHistoryErr := t.DocumentHistoryRepository.GetLastRejection(document.ID.String())
		if fetchHistoryErr != nil {
			return nil, fetchDocErr
		}

		if rejectedBy != nil {
			user, err := t.UserRepository.Get(rejectedBy.UserID.String(), true, orgID)
			if err != nil {
				return nil, err
			}

			positionName := ""
			if user.Position != nil {
				positionName = user.Position.Name
			}

			rejected = response.RejectedOverviewResponse{
				Name:    user.FirstName + " " + user.LastName,
				Title:   positionName,
				Subject: document.Subject,
				Reason:  rejectedBy.Description,
				Date:    rejectedBy.CreatedAt.String(),
			}
		}
	}

	if document != nil {
		return &rejected, nil
	} else {
		return nil, nil
	}
}

func (t DocumentServiceImpl) GetCompletedOverview(userId string, orgID string) (*response.CompletedOverviewResponse, *helper.ErrorModel) {
	document, fetchDocErr := t.DocumentRepository.GetLastestCompleted(userId, orgID)

	completed := response.CompletedOverviewResponse{}

	if fetchDocErr != nil {
		return nil, fetchDocErr
	}

	if document != nil {
		history, fetchHistoryErr := t.DocumentHistoryRepository.GetLastApprover(document.ID.String())
		if fetchHistoryErr != nil {
			return nil, fetchDocErr
		}

		user, err := t.UserRepository.Get(history.UserID.String(), true, orgID)
		if err != nil {
			return nil, err
		}

		recipients, fetchRecipientErr := t.RecipientRepository.GetRecipientsByDocId(document.ID.String())
		if fetchRecipientErr != nil {
			return nil, fetchDocErr
		}

		internalRecipients := []response.InternalRecipientForOverview{}
		if recipients != nil {
			for _, recipient := range recipients {
				user, err := t.UserRepository.Get(recipient.UserID.String(), true, orgID)
				if err != nil {
					return nil, err
				}
				recipientPositionName := ""
				if user.Position != nil {
					recipientPositionName = user.Position.Name
				}

				internalRecipients = append(internalRecipients, response.InternalRecipientForOverview{
					Name:  user.FirstName + " " + user.LastName,
					Title: recipientPositionName,
				})
			}
		}

		positionName := ""
		if user.Position != nil {
			positionName = user.Position.Name
		}

		completed = response.CompletedOverviewResponse{
			IsFinished:        document.Status == 2,
			Name:              user.FirstName + " " + user.LastName,
			Title:             positionName,
			Subject:           document.Subject,
			Date:              history.CreatedAt.String(),
			InternalRecipient: internalRecipients,
			ExternalRecipient: &document.ExternalRecipient,
		}
	}

	if document != nil {
		return &completed, nil
	} else {
		return nil, nil
	}
}

func (t DocumentServiceImpl) GetAllInbox(userId string, orgID string) ([]response.DocumentResponse, *helper.ErrorModel) {
	response, fetchError := t.DocumentRepository.GetAllInbox(userId, orgID)

	if fetchError != nil {
		return nil, fetchError
	}

	return t.mapDocumentsToDocumentResponse(response), nil
}

func (t DocumentServiceImpl) Update(request request.UpdateDocumentRequest, orgID string) (*model.Document, *helper.ErrorModel) {
	errStructure := t.Validate.Struct(request)

	if errStructure != nil {
		msg := "Structure Error"
		return nil, helper.ErrorCatcher(errStructure, 500, &msg)
	}

	trx := t.Db.Begin()
	defer trx.Rollback()

	document, err := t.DocumentRepository.Get(request.Id, orgID)
	if err != nil {
		msg := "Document Not Found"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	document.Type = request.Type
	document.Priority = request.Priority
	document.Subject = request.Subject
	document.Body = request.Body
	document.ExternalRecipient = request.ExternalRecipient
	document.LetterHead = request.LetterHead

	if document.PublicationNumberType == request.PublicationNumberType {
		if request.PublicationNumberType == 3 {
			document.CustomPublicationNumber = request.PublicationValue
		}
	} else {
		document.PublicationNumberType = request.PublicationNumberType
		switch request.PublicationNumberType {
		case 3:
			document.CustomPublicationNumber = request.PublicationValue
		case 4:
			document.CustomPublicationNumber = nil
		}
	}

	if request.IsDraft {
		document.Status = 0
	} else {
		document.Status = 1
		document.Step = 1
	}

	t.DocumentRepository.Update(*document, orgID)
	t.DocumentReferenceRepository.Update(request.References, document.ID)

	// Update Internal Recipients
	if len(request.Recipients) > 0 {
		var recipients []model.Recipient
		for _, userId := range request.Recipients {
			userUuid, _ := uuid.FromString(userId)

			recipients = append(recipients, model.Recipient{
				Document:   document,
				UserID:     userUuid,
				DocumentID: document.ID,
			})
		}

		err := t.RecipientRepository.Update(
			*document,
			recipients,
		)

		if err != nil {
			msg := "Structure Error"
			return nil, helper.ErrorCatcher(err, 500, &msg)
		}
	}

	if len(request.CarbonCopies) > 0 {
		var carbonCopies []model.CarbonCopy
		for _, userId := range request.CarbonCopies {
			userUuid, errorParse := uuid.FromString(userId)
			if errorParse != nil {
				msg := "Structure Error"
				return nil, helper.ErrorCatcher(errorParse, 500, &msg)
			}

			carbonCopies = append(carbonCopies, model.CarbonCopy{
				Document: document,
				UserID:   userUuid,
			})
		}

		err := t.CarbonCopyRepository.Update(
			*document,
			carbonCopies,
		)

		if err != nil {
			msg := "Structure Error"
			return nil, helper.ErrorCatcher(err, 500, &msg)
		}
	}

	if len(request.Sequences) > 0 {
		var sequences []model.DocumentSequence
		for index, sequence := range request.Sequences {
			userUuid, errorParse := uuid.FromString(sequence.UserID)
			if errorParse != nil {
				msg := "Structure Error"
				return nil, helper.ErrorCatcher(errorParse, 500, &msg)
			}

			sequences = append(sequences, model.DocumentSequence{
				DocumentID: &document.ID,
				UserID:     userUuid,
				Step:       index + 1,
				Signature:  sequence.Signature,
			})
		}

		err := t.DocumentSequenceRepository.Update(
			*document,
			sequences,
		)

		if err != nil {
			msg := "Structure Error"
			return nil, helper.ErrorCatcher(err, 500, &msg)
		}
	}

	for _, value := range request.NewAttachments {
		t.DocumentAttachmentRepository.Create(
			trx,
			model.DocumentAttachment{
				Document:     document,
				OriginalName: value.OriginalName,
				FileName:     value.FileName,
				Path:         value.Path,
				Size:         value.Size,
				Type:         value.Type,
			},
		)
	}

	trx.Commit()

	// Notify first approver when document is re-submitted (not a draft)
	if !request.IsDraft && len(request.Sequences) > 0 {
		docID := document.ID.String()
		docSubject := document.Subject
		firstApproverID := request.Sequences[0].UserID
		authorID := request.AuthorID
		frontendURL := t.FrontendURL
		go func() {
			approver, err := t.UserRepository.Get(firstApproverID, true, orgID)
			if err != nil || approver == nil {
				return
			}
			author, err := t.UserRepository.Get(authorID, true, orgID)
			if err != nil || author == nil {
				return
			}
			fromName := author.FirstName + " " + author.LastName
			documentURL := fmt.Sprintf("%s/preview/%s", frontendURL, docID)
			t.EmailService.SendApprovalRequest(approver.Email, approver.FirstName+" "+approver.LastName, fromName, docSubject, documentURL)
		}()
	}

	return document, nil
}

func (t DocumentServiceImpl) Authorize(request request.Authorize, userId string, orgID string) *helper.ErrorModel {
	errStructure := t.Validate.Struct(request)

	if errStructure != nil {
		msg := "Structure Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	document, err := t.DocumentRepository.Get(request.DocumentID, orgID)
	if err != nil {
		msg := "Failed to get document"
		return helper.ErrorCatcher(fmt.Errorf("document not found"), 404, &msg)
	}

	sequences, fetchSequenceErr := t.DocumentSequenceRepository.GetSequencesByDocumentId(document.ID.String())
	if fetchSequenceErr != nil {
		msg := "Failed to get document sequences"
		return helper.ErrorCatcher(fmt.Errorf("sequences not found"), 500, &msg)
	}

	if document == nil || len(sequences) == 0 {
		msg := "Document or sequences not found"
		return helper.ErrorCatcher(fmt.Errorf("not found"), 404, &msg)
	}

	today := time.Now()
	userIdUUID, _ := uuid.FromString(userId)

	// Determine which sequence user corresponds to the current approver.
	// This may be the current user directly, or the current user acting as a delegate.
	var sequenceUserID string
	var onBehalfOfID *uuid.UUID

	for _, seq := range sequences {
		if seq.Step != document.Step {
			continue
		}
		seqUserStr := seq.UserID.String()

		if seqUserStr == userId {
			// Direct authorization
			sequenceUserID = seqUserStr
			break
		}

		// Check if userId is the resolved delegate for this sequence user
		resolved, errResolve := t.resolveDelegate(seqUserStr, today, orgID)
		if errResolve != nil {
			return errResolve
		}
		if resolved == userId {
			sequenceUserID = seqUserStr
			seqUserUUID, _ := uuid.FromString(seqUserStr)
			onBehalfOfID = &seqUserUUID
			break
		}
	}

	if sequenceUserID == "" {
		msg := "You are not authorized to approve this document"
		return helper.ErrorCatcher(fmt.Errorf("unauthorized"), 403, &msg)
	}

	if request.State == 1 { // Approved
		hasSignature := false
		_, errSignature := t.SignatureRepository.GetByUserId(userId, orgID)
		if errSignature == nil {
			hasSignature = true
		}

		for i, seq := range sequences {
			if seq.UserID.String() == sequenceUserID && seq.Step == document.Step {
				sequences[i].Signature = hasSignature
				errUpdateSeq := t.Db.Save(&sequences[i]).Error
				if errUpdateSeq != nil {
					msg := "Failed to update signature status"
					return helper.ErrorCatcher(errUpdateSeq, 500, &msg)
				}
				break
			}
		}

		if (document.Step + 1) <= len(sequences) {
			document.Status = 1
			document.Step = (document.Step + 1)
		} else {
			document.Status = 2
		}
	} else if request.State == 2 { // Rejected
		document.Status = 99
	} else if request.State == 3 { // Cancelled
		document.Status = 3
	}

	isApproved := request.State == 1

	historyEntry := model.DocumentHistory{
		Document:     document,
		Description:  request.Comment,
		UserID:       userIdUUID,
		OnBehalfOfID: onBehalfOfID,
		IsApproved:   isApproved,
	}

	if errResponse := t.DocumentHistoryRepository.Create(historyEntry); errResponse != nil {
		msg := "Failed to create document history"
		return helper.ErrorCatcher(fmt.Errorf("history creation failed"), 500, &msg)
	}

	if errDocumentResponse := t.DocumentRepository.Update(*document, orgID); errDocumentResponse != nil {
		msg := "Failed to update document"
		return helper.ErrorCatcher(fmt.Errorf("document update failed"), 500, &msg)
	}

	// Send email notifications (fire-and-forget, does not affect response)
	docID := document.ID.String()
	docSubject := document.Subject
	docStatus := document.Status
	docStep := document.Step
	authorID := ""
	if document.AuthorID != nil {
		authorID = document.AuthorID.String()
	}
	rejectorID := userId
	rejectorComment := request.Comment
	seqSnapshot := sequences
	frontendURL := t.FrontendURL

	go func() {
		documentURL := fmt.Sprintf("%s/preview/%s", frontendURL, docID)

		switch {
		case request.State == 2 && authorID != "": // Rejected — notify author
			author, err := t.UserRepository.Get(authorID, true, orgID)
			if err != nil || author == nil {
				return
			}
			rejector, err := t.UserRepository.Get(rejectorID, true, orgID)
			if err != nil || rejector == nil {
				return
			}
			rejectorName := rejector.FirstName + " " + rejector.LastName
			t.EmailService.SendDocumentRejected(author.Email, author.FirstName+" "+author.LastName, docSubject, rejectorName, rejectorComment, documentURL)

		case request.State == 1 && docStatus == 2 && authorID != "": // All approved — notify author
			author, err := t.UserRepository.Get(authorID, true, orgID)
			if err != nil || author == nil {
				return
			}
			t.EmailService.SendDocumentApproved(author.Email, author.FirstName+" "+author.LastName, docSubject, documentURL)

		case request.State == 1 && docStatus == 1: // Step advanced — notify next approver
			var nextApproverID string
			for _, seq := range seqSnapshot {
				if seq.Step == docStep {
					nextApproverID = seq.UserID.String()
					break
				}
			}
			if nextApproverID == "" {
				return
			}
			nextApprover, err := t.UserRepository.Get(nextApproverID, true, orgID)
			if err != nil || nextApprover == nil {
				return
			}
			prevApprover, err := t.UserRepository.Get(rejectorID, true, orgID) // rejectorID = current userId (previous approver)
			if err != nil || prevApprover == nil {
				return
			}
			fromName := prevApprover.FirstName + " " + prevApprover.LastName
			t.EmailService.SendApprovalRequest(nextApprover.Email, nextApprover.FirstName+" "+nextApprover.LastName, fromName, docSubject, documentURL)
		}
	}()

	return nil
}

// // User Log
// t.UserLogRepository.Create(
// 	model.UserLog{
// 		UserID: userId,
// 		Action: string(enums.Approve),
// 		Module: string(enums.Document),
// 		Log:    helper.ToJSON(request),
// 	},
// )

func (t DocumentServiceImpl) GetCompleteByAuthorID(authorID string, orgID string) ([]response.DocumentResponse, *helper.ErrorModel) {
	// Ambil data dokumen berdasarkan AuthorID dari repository

	//fmt.Println("Executing query for AuthorID:", authorID)
	documents, fetchError := t.DocumentRepository.GetCompleteByAuthorID(authorID, orgID)
	if fetchError != nil {
		return nil, fetchError
	}

	documentResponses := t.mapDocumentsToDocumentResponse(documents)

	return documentResponses, nil
}

func (t DocumentServiceImpl) GetDraftByAuthorID(authorID string, orgID string) ([]response.DocumentResponse, *helper.ErrorModel) {
	// Ambil data dokumen draft berdasarkan AuthorID dari repository
	documents, fetchError := t.DocumentRepository.GetDraftByAuthorID(authorID, orgID)
	if fetchError != nil {
		return nil, fetchError
	}

	documentResponses := t.mapDocumentsToDocumentResponse(documents)

	return documentResponses, nil
}

func (t DocumentServiceImpl) GetRejectedByAuthorID(authorID string, orgID string) ([]response.DocumentResponse, *helper.ErrorModel) {
	// Ambil data dokumen draft berdasarkan AuthorID dari repository
	documents, fetchError := t.DocumentRepository.GetRejectedByAuthorID(authorID, orgID)
	if fetchError != nil {
		return nil, fetchError
	}

	documentResponses := t.mapDocumentsToDocumentResponse(documents)

	return documentResponses, nil
}

func (t DocumentServiceImpl) GetAllAuthorDocuments(authorID string, orgID string) ([]response.DocumentResponse, *helper.ErrorModel) {
	// Ambil data dokumen draft berdasarkan AuthorID dari repository
	documents, fetchError := t.DocumentRepository.GetAllAuthorDocuments(authorID, orgID)
	if fetchError != nil {
		return nil, fetchError
	}

	documentResponses := t.mapDocumentsToDocumentResponse(documents)

	return documentResponses, nil
}

// Threshold untuk menentukan alert_type.
// Ubah nilai ini sesuai kebutuhan bisnis kantor kamu.
const (
	needApprovalWarningDays = 3 // hari: jika surat pending lebih dari ini → warning
	inProgressWarningDays   = 7 // hari: jika in-progress lebih dari ini → warning
)

func (t DocumentServiceImpl) GetDashboardSummary(userId string, period string, orgID string) (*response.DashboardSummaryResponse, *helper.ErrorModel) {
	raw, fetchErr := t.DocumentRepository.GetDashboardSummary(userId, period, orgID)
	if fetchErr != nil {
		return nil, fetchErr
	}

	// ── NeedApproval card ────────────────────────────────────────────────────
	needApprovalAlertType := "success"
	needApprovalAlertLabel := "Semua berjalan tepat waktu"
	if raw.OldestPendingDays >= needApprovalWarningDays || raw.NeedApprovalUrgent > 0 {
		needApprovalAlertType = "warning"
		needApprovalAlertLabel = fmt.Sprintf("Terlama: %d hari belum diproses", raw.OldestPendingDays)
	}

	// ── InProgress card ──────────────────────────────────────────────────────
	inProgressAlertType := "success"
	inProgressAlertLabel := "Semua berjalan tepat waktu"
	if raw.LongestProcessingDays >= inProgressWarningDays {
		inProgressAlertType = "warning"
		inProgressAlertLabel = fmt.Sprintf("Paling lama: %d hari belum selesai", raw.LongestProcessingDays)
	}

	// ── Rejected card ────────────────────────────────────────────────────────
	rejectedAlertType := "success"
	rejectedAlertLabel := "Tidak ada surat yang ditolak"
	if raw.RejectedTotal > 0 {
		rejectedAlertType = "warning"
		rejectedAlertLabel = fmt.Sprintf("%d surat perlu revisi segera", raw.MineNeedsRevision)
	}

	// ── Completed card ───────────────────────────────────────────────────────
	completedAlertType := "success"
	completedAlertLabel := "Semua berjalan tepat waktu"

	return &response.DashboardSummaryResponse{
		Period: period,
		NeedApproval: response.NeedApprovalCard{
			Total:             raw.NeedApprovalTotal,
			Urgent:            raw.NeedApprovalUrgent,
			Normal:            raw.NeedApprovalNormal,
			OldestPendingDays: raw.OldestPendingDays,
			AlertType:         needApprovalAlertType,
			AlertLabel:        needApprovalAlertLabel,
		},
		InProgress: response.InProgressCard{
			Total:                 raw.InProgressTotal,
			LongestProcessingDays: raw.LongestProcessingDays,
			AlertType:             inProgressAlertType,
			AlertLabel:            inProgressAlertLabel,
		},
		Rejected: response.RejectedCard{
			Total:             raw.RejectedTotal,
			MineNeedsRevision: raw.MineNeedsRevision,
			AlertType:         rejectedAlertType,
			AlertLabel:        rejectedAlertLabel,
		},
		Completed: response.CompletedCard{
			Total:      raw.CompletedTotal,
			TotalYear:  raw.CompletedTotalYear,
			AlertType:  completedAlertType,
			AlertLabel: completedAlertLabel,
		},
	}, nil
}

func (t DocumentServiceImpl) GetDeadlines(userId string, orgID string) ([]response.DeadlineItemResponse, *helper.ErrorModel) {
	documents, err := t.DocumentRepository.GetDeadlines(userId, orgID)
	if err != nil {
		return nil, err
	}

	var result []response.DeadlineItemResponse
	now := time.Now()

	for _, doc := range documents {
		if doc.DueDate == nil {
			continue
		}

		daysRemaining := int(doc.DueDate.Sub(now).Hours() / 24)

		result = append(result, response.DeadlineItemResponse{
			ID:            doc.ID.String(),
			Subject:       doc.Subject,
			DaysRemaining: daysRemaining,
		})
	}

	return result, nil
}

func (t DocumentServiceImpl) GetRecentActivities(userId string, orgID string) ([]response.ActivityResponse, *helper.ErrorModel) {
	histories, err := t.DocumentRepository.GetRecentActivities(userId, orgID)
	if err != nil {
		return nil, err
	}

	var result []response.ActivityResponse

	for _, h := range histories {
		if h.Document == nil {
			continue
		}

		// Ambil nama approver dari UserRepository
		approver, errUser := t.UserRepository.Get(h.UserID.String(), true, orgID)
		if errUser != nil {
			continue
		}

		approverName := approver.FirstName + " " + approver.LastName

		result = append(result, response.ActivityResponse{
			ID:           h.Document.ID.String(),
			Subject:      h.Document.Subject,
			IsApproved:   h.IsApproved,
			ApproverName: approverName,
			UpdatedAt:    h.CreatedAt.Format(time.RFC3339),
		})
	}

	return result, nil
}

func (t DocumentServiceImpl) GetRecentDocuments(userId string, docType int, orgID string) ([]response.RecentDocumentResponse, *helper.ErrorModel) {
	documents, err := t.DocumentRepository.GetRecentDocuments(userId, docType, orgID)
	if err != nil {
		return nil, err
	}

	var result []response.RecentDocumentResponse

	for _, doc := range documents {
		// ── Nomor Surat ──────────────────────────────────────────
		number := "-"
		if doc.CustomPublicationNumber != nil && *doc.CustomPublicationNumber != "" {
			number = *doc.CustomPublicationNumber
		} else {
			docNumber, _ := t.DocumentNumbersRepository.GetByDocumentID(doc.ID, orgID)
			if docNumber != nil {
				number = docNumber.Value
			}
		}

		// ── Status ───────────────────────────────────────────────
		statusMap := map[int]string{
			0:  "Draft",
			1:  "In Progress",
			2:  "Selesai",
			3:  "Cancelled",
			99: "Ditolak",
		}
		statusLabel := statusMap[doc.Status]

		// ── Dari/Kepada ──────────────────────────────────────────
		fromTo := "-"

		if doc.Status == 99 {
			// Rejected: cari approver terakhir dari DocumentHistory
			histories, errHistory := t.DocumentHistoryRepository.GetAllHistoryByDocumentId(doc.ID.String())
			if errHistory == nil && len(histories) > 0 {
				latest := histories[0] // sudah di-sort DESC dari repository
				user, errUser := t.UserRepository.Get(latest.UserID.String(), true, orgID)
				if errUser == nil && user != nil {
					fromTo = user.FirstName + " " + user.LastName
				}
			}
		} else {
			// Status lain: cari dari DocumentSequence berdasarkan step saat ini
			for _, seq := range doc.DocumentSequence {
				if seq.Step == doc.Step {
					user, errUser := t.UserRepository.Get(seq.UserID.String(), true, orgID)
					if errUser == nil && user != nil {
						fromTo = user.FirstName + " " + user.LastName
					}
					break
				}
			}

			// Fallback: jika tidak ada sequence yang match, pakai nama author
			if fromTo == "-" && doc.Author != nil {
				fromTo = doc.Author.FirstName + " " + doc.Author.LastName
			}
		}

		result = append(result, response.RecentDocumentResponse{
			ID:        doc.ID.String(),
			Number:    number,
			Subject:   doc.Subject,
			FromTo:    fromTo,
			Status:    statusLabel,
			Type:      doc.Type,
			UpdatedAt: doc.UpdatedAt.Format(time.RFC3339),
		})
	}

	return result, nil
}

func (s *DocumentServiceImpl) Search(keyword string, orgID string) ([]response.SearchDocumentResponse, *helper.ErrorModel) {
	documents, err := s.DocumentRepository.Search(keyword, orgID)
	if err != nil {
		return nil, err
	}

	var results []response.SearchDocumentResponse
	for _, doc := range documents {
		var docNumber *string

		switch doc.PublicationNumberType {
		case 1, 2:
			docNum, errDocNum := s.DocumentNumbersRepository.GetByDocumentID(doc.ID, orgID)
			if errDocNum == nil && docNum != nil {
				docNumber = &docNum.Value
			}
		case 3:
			docNumber = doc.CustomPublicationNumber
		case 4:
			docNumber = nil
		}

		results = append(results, response.SearchDocumentResponse{
			ID:             doc.ID,
			Subject:        doc.Subject,
			DocumentNumber: docNumber,
		})
	}

	return results, nil
}

func (t DocumentServiceImpl) GetVerification(documentId string) (*response.VerificationResponse, *helper.ErrorModel) {
	// Public, unauthenticated endpoint — no org_id in context. Fetch the
	// document unscoped, then use its own OrganizationID to scope every
	// downstream lookup below.
	document, err := t.DocumentRepository.GetUnscoped(documentId)
	if err != nil || document == nil {
		msg := "Document not found"
		return nil, helper.ErrorCatcher(fmt.Errorf("not found"), 404, &msg)
	}
	orgID := document.OrganizationID.String()

	if document.Status != 2 {
		msg := "Document is not yet finalized"
		return nil, helper.ErrorCatcher(fmt.Errorf("not finished"), 404, &msg)
	}

	// Document number: prefer custom, else from document_numbers table
	docNumber := ""
	switch document.PublicationNumberType {
	case 3:
		if document.CustomPublicationNumber != nil {
			docNumber = *document.CustomPublicationNumber
		}
	case 1, 2:
		docNum, errDocNum := t.DocumentNumbersRepository.GetByDocumentID(document.ID, orgID)
		if errDocNum == nil && docNum != nil {
			docNumber = docNum.Value
		}
	}

	// Fetch all app_settings at once for efficiency
	settingsMap := map[string]string{}
	if allSettings, errSettings := t.AppSettingsRepository.GetAll(orgID); errSettings == nil {
		for _, s := range allSettings {
			settingsMap[s.Key] = s.Value
		}
	}
	get := func(key string) string { return settingsMap[key] }

	// Last approver name (last history record)
	lastApprover, _ := t.DocumentHistoryRepository.GetLastApprover(documentId)
	lastApproverName := ""
	if lastApprover != nil {
		actorID := lastApprover.UserID.String()
		if lastApprover.OnBehalfOfID != nil {
			actorID = lastApprover.OnBehalfOfID.String()
		}
		user, errUser := t.UserRepository.Get(actorID, true, orgID)
		if errUser == nil && user != nil {
			lastApproverName = user.FirstName + " " + user.LastName
		}
	}

	// Approval date from document.UpdatedAt
	approvalDate := ""
	if document.UpdatedAt != nil {
		approvalDate = document.UpdatedAt.Format("02 January 2006")
	}

	return &response.VerificationResponse{
		Subject:            document.Subject,
		Body:               document.Body,
		DocumentNumber:     docNumber,
		OrganizationName:   get("company_name"),
		ApprovalDate:       approvalDate,
		LastApproverName:   lastApproverName,
		Type:               document.Type,
		CompanyLogoUrl:     get("company_logo_url"),
		CompanyAddress:     get("company_address"),
		CompanyCity:        get("company_city"),
		CompanyPhone:       get("company_phone_number"),
		CompanyEmail:       get("company_email"),
		CompanyDescription: get("company_description"),
	}, nil
}

func (t DocumentServiceImpl) Recall(documentId string, userId string, orgID string) *helper.ErrorModel {
	document, err := t.DocumentRepository.Get(documentId, orgID)
	if err != nil {
		msg := "Document not found"
		return helper.ErrorCatcher(fmt.Errorf("not found"), 404, &msg)
	}

	if document.Author.ID.String() != userId {
		msg := "Only the document author can recall this document"
		return helper.ErrorCatcher(fmt.Errorf("forbidden"), 403, &msg)
	}

	if document.Status != 1 {
		msg := "Only in-progress documents can be recalled"
		return helper.ErrorCatcher(fmt.Errorf("invalid status"), 400, &msg)
	}

	histories, errHistory := t.DocumentHistoryRepository.GetAllHistoryByDocumentId(documentId)
	if errHistory != nil {
		msg := "Failed to check document history"
		return helper.ErrorCatcher(fmt.Errorf("history check failed"), 500, &msg)
	}
	if len(histories) > 0 {
		msg := "Document cannot be recalled after an approver has already acted"
		return helper.ErrorCatcher(fmt.Errorf("already actioned"), 400, &msg)
	}

	document.Status = 0
	document.Step = 1
	if errUpdate := t.DocumentRepository.Update(*document, orgID); errUpdate != nil {
		msg := "Failed to recall document"
		return helper.ErrorCatcher(fmt.Errorf("update failed"), 500, &msg)
	}

	userUUID, _ := uuid.FromString(userId)
	t.UserLogRepository.Create(model.UserLog{
		UserID: userUUID,
		Action: "Recall",
		Module: "Document",
	}, orgID)

	return nil
}
