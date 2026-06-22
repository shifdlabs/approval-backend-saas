package document

type DashboardSummaryResponse struct {
	Period       string           `json:"period"`
	NeedApproval NeedApprovalCard `json:"need_approval"`
	InProgress   InProgressCard   `json:"in_progress"`
	Rejected     RejectedCard     `json:"rejected"`
	Completed    CompletedCard    `json:"completed"`
}

type NeedApprovalCard struct {
	Total             int    `json:"total"`
	Urgent            int    `json:"urgent"`
	Normal            int    `json:"normal"`
	OldestPendingDays int    `json:"oldest_pending_days"`
	AlertType         string `json:"alert_type"` // "warning" | "success"
	AlertLabel        string `json:"alert_label"`
}

type InProgressCard struct {
	Total                 int    `json:"total"`
	LongestProcessingDays int    `json:"longest_processing_days"`
	AlertType             string `json:"alert_type"`
	AlertLabel            string `json:"alert_label"`
}

type RejectedCard struct {
	Total             int    `json:"total"`
	MineNeedsRevision int    `json:"mine_needs_revision"`
	AlertType         string `json:"alert_type"`
	AlertLabel        string `json:"alert_label"`
}

type CompletedCard struct {
	Total      int    `json:"total"`
	TotalYear  int    `json:"total_year"`
	AlertType  string `json:"alert_type"`
	AlertLabel string `json:"alert_label"`
}
