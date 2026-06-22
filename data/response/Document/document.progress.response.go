package document

type DashboardProgressResponse struct {
	InProgress *DocumentInProgressResponse `json:"inProgress"`
	Rejected   *RejectedOverviewResponse   `json:"rejected"`
	Completed  *CompletedOverviewResponse  `json:"completed"`
}

type DocumentInProgressResponse struct {
	Subject   string                `json:"subject"`
	Approvers []ApproverForOverview `json:"approvers"`
}

type ApproverForOverview struct {
	Name         string  `json:"name"`
	Title        string  `json:"title"`
	Approved     *bool   `json:"approved"`
	Date         *string `json:"date"`
	Signature    bool    `json:"signature"`
	SignatureUrl *string `json:"signatureUrl"`
	DelegateName *string `json:"delegateName"` // set when pending step has active delegation
	OnBehalfOf   *string `json:"onBehalfOf"`   // set when approved by a delegate; value = delegate's name
}

type RejectedOverviewResponse struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	Subject string `json:"subject"`
	Reason  string `json:"reason"`
	Date    string `json:"date"`
}

type CompletedOverviewResponse struct {
	IsFinished        bool                           `json:"isFinished"`
	Name              string                         `json:"name"`
	Title             string                         `json:"title"`
	Subject           string                         `json:"subject"`
	Date              string                         `json:"date"`
	InternalRecipient []InternalRecipientForOverview `json:"internalRecipient"`
	ExternalRecipient *string                        `json:"externalRecipient"`
}

type InternalRecipientForOverview struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}
