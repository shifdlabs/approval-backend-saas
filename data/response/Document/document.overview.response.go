package document

type DocumentOverviewResponse struct {
	InProgress int `json:"inProgress"`
	Finished   int `json:"finished"`
	Cancelled  int `json:"cancelled"`
	TotalValue int `json:"totalValue"`
}
