package document

type DashboardResponse struct {
	Statistic       DocumentStatisticResponse `json:"statistic"`
	AuthorDocuments []DocumentResponse        `json:"authorDocuments"`
	Progress        DashboardProgressResponse `json:"progress"`
	Inbox           []DocumentResponse        `json:"inbox"`
}
