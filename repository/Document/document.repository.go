package document

import (
	dashboard "Microservice/data/model/Dashboard"
	"Microservice/helper"
	"Microservice/model"

	"gorm.io/gorm"
)

type DocumentRepository interface {
	Create(db gorm.DB, report *model.Document) *helper.ErrorModel
	Get(id string, orgID string) (*model.Document, *helper.ErrorModel)
	GetAll(orgID string) ([]model.Document, *helper.ErrorModel)
	GetAllReferences(query string, orgID string) ([]model.Document, *helper.ErrorModel)
	GetAllAuthorization(id string, orgID string) ([]model.Document, *helper.ErrorModel)
	GetAllInbox(id string, orgID string) ([]model.Document, *helper.ErrorModel)
	GetAllInProgress(userId string, orgID string) ([]model.Document, *helper.ErrorModel)
	GetDocumentStatistics(id string, orgID string) ([]int, *helper.ErrorModel)
	GetOneLatestInprogress(id string, orgID string) (*model.Document, *helper.ErrorModel)
	GetLastestRejected(id string, orgID string) (*model.Document, *helper.ErrorModel)
	GetLastestCompleted(id string, orgID string) (*model.Document, *helper.ErrorModel)
	Update(report model.Document, orgID string) *helper.ErrorModel
	Delete(id string, orgID string) *helper.ErrorModel
	GetCompleteByAuthorID(authorID string, orgID string) ([]model.Document, *helper.ErrorModel)
	GetDraftByAuthorID(authorID string, orgID string) ([]model.Document, *helper.ErrorModel)
	GetRejectedByAuthorID(authorID string, orgID string) ([]model.Document, *helper.ErrorModel)
	GetAllAuthorDocuments(authorID string, orgID string) ([]model.Document, *helper.ErrorModel)
	GetDashboardSummary(userId string, period string, orgID string) (*dashboard.DashboardSummaryRaw, *helper.ErrorModel)
	GetDeadlines(userId string, orgID string) ([]model.Document, *helper.ErrorModel)
	GetRecentActivities(userId string, orgID string) ([]model.DocumentHistory, *helper.ErrorModel)
	GetRecentDocuments(userId string, docType int, orgID string) ([]model.Document, *helper.ErrorModel)
	Search(keyword string, orgID string) ([]model.Document, *helper.ErrorModel)
	GetAllInProgressForSLA() ([]model.Document, *helper.ErrorModel)

	// GetUnscoped bypasses org filtering. Only for the public, unauthenticated
	// verification endpoint (GET /api/verification/:id), which has no org_id in
	// context. Once fetched, the returned document's own OrganizationID should
	// be used to scope every subsequent lookup in that flow.
	GetUnscoped(id string) (*model.Document, *helper.ErrorModel)
}
