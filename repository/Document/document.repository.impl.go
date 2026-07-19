package document

import (
	dashboard "Microservice/data/model/Dashboard"
	"Microservice/helper"
	"Microservice/model"
	"errors"
	"strings"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type DocumentRepositoryImpl struct {
	Db *gorm.DB
}

func NewDocumentRepositoryImpl(Db *gorm.DB) DocumentRepository {
	return &DocumentRepositoryImpl{Db: Db}
}

func (t *DocumentRepositoryImpl) Create(db gorm.DB, document *model.Document) *helper.ErrorModel {
	result := db.Create(document)
	if result.Error != nil {
		msg := "Create Document Failed"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

func (t *DocumentRepositoryImpl) Get(id string, orgID string) (*model.Document, *helper.ErrorModel) {
	var report model.Document

	reportId, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse id"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.Preload("Author").Preload("DocumentAttachment").Preload("DocumentSequence").Preload("DocumentHistory").Where("organization_id = ?", orgID).First(&report, "id = ?", reportId)

	if result.Error != nil {
		msg := "Get Document Failed"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return &report, nil
}

// GetUnscoped looks up a document by id with no organization filter. Only for
// the public verification endpoint (see interface doc comment).
func (t *DocumentRepositoryImpl) GetUnscoped(id string) (*model.Document, *helper.ErrorModel) {
	var report model.Document

	reportId, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse id"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.Preload("Author").Preload("DocumentAttachment").Preload("DocumentSequence").Preload("DocumentHistory").First(&report, "id = ?", reportId)

	if result.Error != nil {
		msg := "Get Document Failed"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return &report, nil
}

func (t *DocumentRepositoryImpl) GetAll(orgID string) ([]model.Document, *helper.ErrorModel) {
	var reports []model.Document
	result := t.Db.Preload("Author").Preload("DocumentAttachment").Preload("DocumentSequence").Preload("DocumentHistory").Preload("DocumentHistory").Where("organization_id = ?", orgID).Find(&reports)
	if result.Error != nil {
		msg := "Failed to get all documents"
		return reports, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return reports, nil
}

func (t *DocumentRepositoryImpl) GetAllReferences(query string, orgID string) ([]model.Document, *helper.ErrorModel) {
	var reports []model.Document
	result := t.Db.Where("organization_id = ? AND status = ? AND subject ILIKE ?", orgID, 2, "%"+query+"%").Find(&reports)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "Failed to get references"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return reports, nil
}

func (t *DocumentRepositoryImpl) GetAllAuthorization(id string, orgID string) ([]model.Document, *helper.ErrorModel) {
	var reports []model.Document
	result := t.Db.
		Preload("Author").
		Preload("DocumentHistory").
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("documents.organization_id = ?", orgID).
		Where("document_sequences.user_id = ?", id).
		Where("document_sequences.step = documents.step").
		Where("documents.status = 1").
		Find(&reports)

	if result.Error != nil {
		msg := "Failed to get all documents"
		return reports, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return reports, nil
}

func (t *DocumentRepositoryImpl) GetDocumentStatistics(id string, orgID string) ([]int, *helper.ErrorModel) {
	var totalAuthorization int64

	var totalInProgressAsApprover int64
	var totalInProgressAsAuthor int64

	var totalRejectedAsApprover int64
	var totalRejectedAsAuthor int64

	var totalCompletedAsApprover int64
	var totalCompletedAsAuthor int64

	countAuthorization := t.Db.
		Model(&model.Document{}).
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("documents.organization_id = ?", orgID).
		Where("document_sequences.user_id = ?", id).
		Where("document_sequences.step = documents.step").
		Count(&totalAuthorization)

	if countAuthorization.Error != nil {
		msg := "Failed to get all documents"
		return nil, helper.ErrorCatcher(countAuthorization.Error, 500, &msg)
	}

	countInProgressAsApprover := t.Db.
		Model(&model.Document{}).
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("documents.organization_id = ?", orgID).
		Where("document_sequences.user_id = ?", id).
		Where("documents.status = 1 OR documents.status = 99").
		Count(&totalInProgressAsApprover)

	if countInProgressAsApprover.Error != nil {
		msg := "Failed to get all documents"
		return nil, helper.ErrorCatcher(countInProgressAsApprover.Error, 500, &msg)
	}

	countInProgressAsAuthor := t.Db.
		Model(&model.Document{}).
		Where("organization_id = ?", orgID).
		Where("status = 1").
		Where("author_id = ?", id).
		Count(&totalInProgressAsAuthor)

	if countInProgressAsAuthor.Error != nil {
		msg := "Failed to get all documents"
		return nil, helper.ErrorCatcher(countInProgressAsAuthor.Error, 500, &msg)
	}

	countRejectedAsApprover := t.Db.
		Model(&model.Document{}).
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("documents.organization_id = ?", orgID).
		Where("document_sequences.user_id = ?", id).
		Where("documents.status = 99").
		Count(&totalRejectedAsApprover)

	if countRejectedAsApprover.Error != nil {
		msg := "Failed to get all documents"
		return nil, helper.ErrorCatcher(countRejectedAsApprover.Error, 500, &msg)
	}

	countRejectedAsAuthor := t.Db.
		Model(&model.Document{}).
		Where("organization_id = ?", orgID).
		Where("status = 99").
		Where("author_id = ?", id).
		Count(&totalRejectedAsAuthor)

	if countRejectedAsAuthor.Error != nil {
		msg := "Failed to get all documents"
		return nil, helper.ErrorCatcher(countRejectedAsAuthor.Error, 500, &msg)
	}

	countCompletedAsApprover := t.Db.
		Model(&model.Document{}).
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("documents.organization_id = ?", orgID).
		Where("document_sequences.user_id = ?", id).
		Where("documents.status = 2 OR documents.status = 3").
		Count(&totalCompletedAsApprover)

	if countCompletedAsApprover.Error != nil {
		msg := "Failed to get all documents"
		return nil, helper.ErrorCatcher(countCompletedAsApprover.Error, 500, &msg)
	}

	countCompletedAsAuthor := t.Db.
		Model(&model.Document{}).
		Where("organization_id = ?", orgID).
		Where("status = 2 OR status = 3").
		Where("author_id = ?", id).
		Count(&totalCompletedAsAuthor)

	if countCompletedAsAuthor.Error != nil {
		msg := "Failed to get all documents"
		return nil, helper.ErrorCatcher(countCompletedAsAuthor.Error, 500, &msg)
	}

	var result = []int{
		int(totalAuthorization),
		int(totalInProgressAsApprover + totalInProgressAsAuthor),
		int(totalRejectedAsApprover + totalRejectedAsAuthor),
		int(totalCompletedAsApprover + totalCompletedAsAuthor),
	}

	return result, nil
}

func (t *DocumentRepositoryImpl) GetOneLatestInprogress(id string, orgID string) (*model.Document, *helper.ErrorModel) {
	var doc model.Document

	response := t.Db.
		Model(&model.Document{}).
		Where("organization_id = ?", orgID).
		Where("status = 1").
		Where("author_id = ?", id).
		Order("created_at DESC").
		First(&doc)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "Failed to get in-progress document"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return &doc, nil
}

func (t *DocumentRepositoryImpl) GetLastestRejected(id string, orgID string) (*model.Document, *helper.ErrorModel) {
	var doc model.Document

	response := t.Db.
		Model(&model.Document{}).
		Where("organization_id = ?", orgID).
		Where("status = 99").
		Where("author_id = ?", id).
		Order("created_at DESC").
		First(&doc)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "Failed to get in-progress document"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return &doc, nil
}

func (t *DocumentRepositoryImpl) GetLastestCompleted(id string, orgID string) (*model.Document, *helper.ErrorModel) {
	var doc model.Document

	response := t.Db.
		Model(&model.Document{}).
		Where("organization_id = ?", orgID).
		Where("status IN ?", []int{2, 3}).
		Where("author_id = ?", id).
		Order("created_at DESC").
		First(&doc)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "Failed to get in-progress document"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return &doc, nil
}

func (t *DocumentRepositoryImpl) GetAllInbox(id string, orgID string) ([]model.Document, *helper.ErrorModel) {
	var documents []model.Document

	response := t.Db.
		Model(&model.Document{}).
		Preload("Author"). // Preload relasi Author jika diperlukan
		Joins("JOIN recipients ON recipients.document_id = documents.id").
		Where("documents.organization_id = ?", orgID).
		Where("recipients.user_id = ?", id).
		Where("documents.status = 2").
		Find(&documents)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			// Kembalikan slice kosong jika tidak ada data
			return []model.Document{}, nil
		}

		msg := "Failed to get all documents"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return documents, nil
}

func (t *DocumentRepositoryImpl) GetAllInProgress(userId string, orgID string) ([]model.Document, *helper.ErrorModel) {
	var documents []model.Document

	response := t.Db.
		Model(&model.Document{}).
		Select("DISTINCT documents.*").
		Preload("Author"). // Preload relasi Author jika diperlukan
		Preload("DocumentSequence").
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("documents.organization_id = ?", orgID).
		Where("documents.author_id = ? AND documents.status = 1", userId). // first condition
		Or("document_sequences.user_id = ? AND document_sequences.step > documents.step AND documents.status = 1 AND documents.organization_id = ?", userId, orgID).
		Find(&documents)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			// Kembalikan slice kosong jika tidak ada data
			return []model.Document{}, nil
		}

		msg := "Failed to get all documents"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return documents, nil
}

func (t *DocumentRepositoryImpl) Update(report model.Document, orgID string) *helper.ErrorModel {
	var existing model.Document
	if err := t.Db.Where("organization_id = ?", orgID).First(&existing, "id = ?", report.ID).Error; err != nil {
		msg := "Document not found"
		return helper.ErrorCatcher(err, 404, &msg)
	}

	err := t.Db.Save(&report).Error

	if err != nil {
		msg := "Failed to update document"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	return nil
}

func (t *DocumentRepositoryImpl) Delete(id string, orgID string) *helper.ErrorModel {
	reportId, err := uuid.FromString(id)
	if err != nil {
		msg := "Failed to parse id"
		return helper.ErrorCatcher(err, 500, &msg)
	}

	result := t.Db.Unscoped().Where("organization_id = ?", orgID).Delete(&model.Document{}, reportId)

	if result.Error != nil {
		msg := "Failed to delete document"
		return helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return nil
}

func (t *DocumentRepositoryImpl) GetCompleteByAuthorID(authorID string, orgID string) ([]model.Document, *helper.ErrorModel) {
	var documents []model.Document

	// Gunakan Preload untuk memuat relasi Author dan tambahkan filter status = 2
	result := t.Db.Preload("Author").
		Where("organization_id = ? AND author_id = ? AND status = ?", orgID, authorID, 2).
		Order("updated_at DESC").
		Find(&documents)
	if result.Error != nil {
		msg := "Failed to fetch documents for the author"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return documents, nil
}

func (t *DocumentRepositoryImpl) GetDraftByAuthorID(authorID string, orgID string) ([]model.Document, *helper.ErrorModel) {
	var documents []model.Document

	// Gunakan Preload untuk memuat relasi Author dan tambahkan filter status = 0
	result := t.Db.Preload("Author").
		Where("organization_id = ? AND author_id = ? AND status = ?", orgID, authorID, 0).
		Order("updated_at DESC").
		Find(&documents)
	if result.Error != nil {
		msg := "Failed to fetch draft documents for the author"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return documents, nil
}

func (t *DocumentRepositoryImpl) GetRejectedByAuthorID(authorID string, orgID string) ([]model.Document, *helper.ErrorModel) {
	var documents []model.Document

	// Gunakan Preload untuk memuat relasi Author dan tambahkan filter status = 0
	response := t.Db.Preload("Author").
		Where("organization_id = ? AND author_id = ? AND status = ?", orgID, authorID, 99).
		Order("updated_at DESC").
		Find(&documents)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "Failed to get in-progress document"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return documents, nil
}

func (t *DocumentRepositoryImpl) GetAllAuthorDocuments(authorID string, orgID string) ([]model.Document, *helper.ErrorModel) {
	var documents []model.Document

	// Gunakan Preload untuk memuat relasi Author dan tambahkan filter status = 0
	response := t.Db.Preload("Author").
		Where("organization_id = ? AND author_id = ?", orgID, authorID).
		Order("updated_at DESC").
		Find(&documents)

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		msg := "Failed to get in-progress document"
		return nil, helper.ErrorCatcher(response.Error, 500, &msg)
	}

	return documents, nil
}

// buildPeriodFilter menghasilkan kondisi SQL berdasarkan period yang dipilih.
// Filter diterapkan pada kolom documents.updated_at.
func buildPeriodFilter(period string) string {
	switch period {
	case "today":
		return "documents.updated_at >= CURRENT_DATE"
	case "week":
		return "documents.updated_at >= DATE_TRUNC('week', CURRENT_DATE)"
	case "month":
		return "documents.updated_at >= DATE_TRUNC('month', CURRENT_DATE)"
	default:
		return "1=1" // "all": tidak ada filter tanggal
	}
}

// Struct lokal untuk menampung hasil agregasi GORM Scan
type ageResult struct {
	Days int `gorm:"column:days"`
}

func (t *DocumentRepositoryImpl) GetDashboardSummary(userId string, period string, orgID string) (*dashboard.DashboardSummaryRaw, *helper.ErrorModel) {
	periodFilter := buildPeriodFilter(period)
	raw := &dashboard.DashboardSummaryRaw{}

	// ── 1. NEED APPROVAL ─────────────────────────────────────────────────────
	// Dokumen status=1 yang sudah giliran user ini di sequence (step match).

	var needApprovalTotal int64
	err := t.Db.Model(&model.Document{}).
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("documents.organization_id = ?", orgID).
		Where("document_sequences.user_id = ?", userId).
		Where("document_sequences.step = documents.step").
		Where("documents.status = 1").
		Where(periodFilter).
		Count(&needApprovalTotal)
	if err.Error != nil {
		msg := "Failed to count need_approval total"
		return nil, helper.ErrorCatcher(err.Error, 500, &msg)
	}
	raw.NeedApprovalTotal = int(needApprovalTotal)

	var needApprovalUrgent int64
	t.Db.Model(&model.Document{}).
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("documents.organization_id = ?", orgID).
		Where("document_sequences.user_id = ?", userId).
		Where("document_sequences.step = documents.step").
		Where("documents.status = 1").
		Where("documents.priority = 1"). // Priority 1 = High
		Where(periodFilter).
		Count(&needApprovalUrgent)
	raw.NeedApprovalUrgent = int(needApprovalUrgent)
	raw.NeedApprovalNormal = raw.NeedApprovalTotal - raw.NeedApprovalUrgent

	// oldest_pending_days: dari SELURUH data (tidak ikut filter period),
	// supaya user selalu tau surat tertua yang menunggu tanda tangannya.
	var oldestAge ageResult
	t.Db.Model(&model.Document{}).
		Select("COALESCE(EXTRACT(DAY FROM NOW() - MIN(documents.created_at))::int, 0) as days").
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("documents.organization_id = ?", orgID).
		Where("document_sequences.user_id = ?", userId).
		Where("document_sequences.step = documents.step").
		Where("documents.status = 1").
		Scan(&oldestAge)
	raw.OldestPendingDays = oldestAge.Days

	// ── 2. IN PROGRESS ───────────────────────────────────────────────────────
	// Dokumen status=1 dimana user adalah author ATAU ada di sequence.

	var inProgressTotal int64
	t.Db.Model(&model.Document{}).
		Select("COUNT(DISTINCT documents.id)").
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("documents.organization_id = ?", orgID).
		Where("documents.status = 1").
		Where(periodFilter).
		Where("documents.author_id = ? OR document_sequences.user_id = ?", userId, userId).
		Count(&inProgressTotal)
	raw.InProgressTotal = int(inProgressTotal)

	// longest_processing_days: dari SELURUH data (tidak ikut filter period)
	var longestAge ageResult
	t.Db.Model(&model.Document{}).
		Select("COALESCE(EXTRACT(DAY FROM NOW() - MIN(documents.created_at))::int, 0) as days").
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("documents.organization_id = ?", orgID).
		Where("documents.status = 1").
		Where("documents.author_id = ? OR document_sequences.user_id = ?", userId, userId).
		Scan(&longestAge)
	raw.LongestProcessingDays = longestAge.Days

	// ── 3. REJECTED ──────────────────────────────────────────────────────────
	// Dokumen status=99 dimana user adalah author ATAU ada di sequence.

	var rejectedTotal int64
	t.Db.Model(&model.Document{}).
		Select("COUNT(DISTINCT documents.id)").
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("documents.organization_id = ?", orgID).
		Where("documents.status = 99").
		Where(periodFilter).
		Where("documents.author_id = ? OR document_sequences.user_id = ?", userId, userId).
		Count(&rejectedTotal)
	raw.RejectedTotal = int(rejectedTotal)

	// mine_needs_revision: hanya surat MILIK user yang ditolak
	var mineNeedsRevision int64
	t.Db.Model(&model.Document{}).
		Where("organization_id = ?", orgID).
		Where("status = 99").
		Where("author_id = ?", userId).
		Count(&mineNeedsRevision)
	raw.MineNeedsRevision = int(mineNeedsRevision)

	// ── 4. COMPLETED ─────────────────────────────────────────────────────────
	// Dokumen status=2 atau 3 (finished / cancelled).

	var completedTotal int64
	t.Db.Model(&model.Document{}).
		Select("COUNT(DISTINCT documents.id)").
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("documents.organization_id = ?", orgID).
		Where("documents.status = 2 OR documents.status = 3").
		Where(periodFilter).
		Where("documents.author_id = ? OR document_sequences.user_id = ?", userId, userId).
		Count(&completedTotal)
	raw.CompletedTotal = int(completedTotal)

	// total_year: SELALU aggregate per tahun berjalan, tidak ikut filter period
	var completedTotalYear int64
	t.Db.Model(&model.Document{}).
		Select("COUNT(DISTINCT documents.id)").
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("documents.organization_id = ?", orgID).
		Where("documents.status = 2 OR documents.status = 3").
		Where("EXTRACT(YEAR FROM documents.updated_at) = EXTRACT(YEAR FROM NOW())").
		Where("documents.author_id = ? OR document_sequences.user_id = ?", userId, userId).
		Count(&completedTotalYear)
	raw.CompletedTotalYear = int(completedTotalYear)

	return raw, nil
}

func (t *DocumentRepositoryImpl) GetDeadlines(userId string, orgID string) ([]model.Document, *helper.ErrorModel) {
	var documents []model.Document

	result := t.Db.
		Model(&model.Document{}).
		Select("DISTINCT documents.*").
		Joins("JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("documents.organization_id = ?", orgID).
		Where("documents.due_date IS NOT NULL").
		Where("documents.status = 1"). // hanya yang in-progress
		Where("documents.author_id = ? OR document_sequences.user_id = ?", userId, userId).
		Order("documents.due_date ASC").
		Limit(5).
		Find(&documents)

	if result.Error != nil {
		msg := "Failed to get deadlines"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return documents, nil
}

func (t *DocumentRepositoryImpl) GetRecentActivities(userId string, orgID string) ([]model.DocumentHistory, *helper.ErrorModel) {
	var histories []model.DocumentHistory

	result := t.Db.
		Preload("Document").
		Joins("JOIN documents ON documents.id = document_histories.document_id").
		Where("documents.organization_id = ?", orgID).
		Where("documents.author_id = ?", userId).
		Order("document_histories.created_at DESC").
		Limit(5).
		Find(&histories)

	if result.Error != nil {
		msg := "Failed to get recent activities"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return histories, nil
}

func (t *DocumentRepositoryImpl) GetRecentDocuments(userId string, docType int, orgID string) ([]model.Document, *helper.ErrorModel) {
	var documents []model.Document

	query := t.Db.
		Model(&model.Document{}).
		Select("DISTINCT documents.*").
		Preload("Author").
		Preload("DocumentSequence").
		Preload("DocumentHistory").
		Joins("LEFT JOIN document_sequences ON document_sequences.document_id = documents.id").
		Where("documents.organization_id = ?", orgID).
		Where("documents.author_id = ? OR document_sequences.user_id = ?", userId, userId).
		Where("documents.status != 0"). // exclude draft
		Order("documents.updated_at DESC").
		Limit(10)

	if docType == 1 || docType == 2 {
		query = query.Where("documents.type = ?", docType)
	}

	result := query.Find(&documents)
	if result.Error != nil {
		msg := "Failed to get recent documents"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}

	return documents, nil
}

func (r *DocumentRepositoryImpl) GetAllInProgressForSLA() ([]model.Document, *helper.ErrorModel) {
	var documents []model.Document
	result := r.Db.
		Preload("DocumentSequence").
		Where("status = 1").
		Find(&documents)
	if result.Error != nil {
		msg := "Failed to get in-progress documents for SLA"
		return nil, helper.ErrorCatcher(result.Error, 500, &msg)
	}
	return documents, nil
}

func (r *DocumentRepositoryImpl) Search(keyword string, orgID string) ([]model.Document, *helper.ErrorModel) {
	var documents []model.Document

	// Lowercase keyword sebelum dipakai
	likeKeyword := "%" + strings.ToLower(keyword) + "%"

	err := r.Db.
		Joins("LEFT JOIN document_numbers dn ON dn.document_id = documents.id AND dn.deleted_at IS NULL").
		Where("documents.organization_id = ?", orgID).
		Where(`
            documents.deleted_at IS NULL
            AND (
                LOWER(documents.subject) LIKE ?
                OR (documents.publication_number_type IN (1,2) AND LOWER(dn.value) LIKE ?)
                OR (documents.publication_number_type = 3 AND LOWER(documents.custom_publication_number) LIKE ?)
            )
        `, likeKeyword, likeKeyword, likeKeyword).
		Find(&documents).Error

	if err != nil {
		return nil, &helper.ErrorModel{Code: 500, Message: err.Error()}
	}
	return documents, nil
}
