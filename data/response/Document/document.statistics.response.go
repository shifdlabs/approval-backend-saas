package document

type DocumentStatisticResponse struct {
	Authorization int `json:"authorization"`
	InProgress    int `json:"inProgress"`
	Rejected      int `json:"rejected"`
	Completed     int `json:"completed"`
}
