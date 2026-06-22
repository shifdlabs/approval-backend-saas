package document

import (
	request "Microservice/data/request/Document"
	response "Microservice/data/response/Document"
	"Microservice/helper"
	"Microservice/model"
)

type DocumentService interface {
	Create(request request.CreateDocumentRequest, orgID string) (*model.Document, *helper.ErrorModel)
	GetDocument(id string, orgID string) (*response.DocumentResponse, *helper.ErrorModel)
	GetDetailDocument(id string, currentUserId string, orgID string) (*response.DocumentDetailResponse, *helper.ErrorModel)
	GetDetailForEdit(id string, orgID string) (*response.EditDocumentResponse, *helper.ErrorModel)
	GetAllDocument(orgID string) ([]response.DocumentResponse, *helper.ErrorModel)
	GetAllReferences(query string, orgID string) ([]response.DocumentResponse, *helper.ErrorModel)
	GetAllAuthorization(userId string, orgID string) ([]response.DocumentResponse, *helper.ErrorModel)
	GetAllInbox(userId string, orgID string) ([]response.DocumentResponse, *helper.ErrorModel)
	GetAllInProgress(userId string, orgID string) ([]response.DocumentResponse, *helper.ErrorModel)
	GetDocumentStatistics(userId string, orgID string) (*response.DocumentStatisticResponse, *helper.ErrorModel)
	GetInProgressOverview(userId string, orgID string) (*response.DocumentInProgressResponse, *helper.ErrorModel)
	GetInProgressOverviewByDocId(documentId string, orgID string) (*response.DocumentInProgressResponse, *helper.ErrorModel)
	GetRejectedOverview(userId string, orgID string) (*response.RejectedOverviewResponse, *helper.ErrorModel)
	GetCompletedOverview(userId string, orgID string) (*response.CompletedOverviewResponse, *helper.ErrorModel)
	Update(request request.UpdateDocumentRequest, orgID string) (*model.Document, *helper.ErrorModel)
	Authorize(request request.Authorize, userId string, orgID string) *helper.ErrorModel
	GetCompleteByAuthorID(authorID string, orgID string) ([]response.DocumentResponse, *helper.ErrorModel)
	GetDraftByAuthorID(authorID string, orgID string) ([]response.DocumentResponse, *helper.ErrorModel)
	GetRejectedByAuthorID(authorID string, orgID string) ([]response.DocumentResponse, *helper.ErrorModel)
	GetAllAuthorDocuments(authorID string, orgID string) ([]response.DocumentResponse, *helper.ErrorModel)
	GetDashboardSummary(userId string, period string, orgID string) (*response.DashboardSummaryResponse, *helper.ErrorModel)
	GetDeadlines(userId string, orgID string) ([]response.DeadlineItemResponse, *helper.ErrorModel)
	GetRecentActivities(userId string, orgID string) ([]response.ActivityResponse, *helper.ErrorModel)
	GetRecentDocuments(userId string, docType int, orgID string) ([]response.RecentDocumentResponse, *helper.ErrorModel)
	Search(keyword string, orgID string) ([]response.SearchDocumentResponse, *helper.ErrorModel)
	Recall(documentId string, userId string, orgID string) *helper.ErrorModel

	// GetVerification is the sole exception — it backs the public,
	// unauthenticated GET /api/verification/:id endpoint, so it takes no
	// orgID. It resolves the document unscoped, then scopes every downstream
	// lookup using that document's own OrganizationID.
	GetVerification(documentId string) (*response.VerificationResponse, *helper.ErrorModel)
}
