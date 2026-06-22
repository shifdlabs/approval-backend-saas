package delegator

type UpdateDelegatorRequest struct {
	DelegateID string `validate:"required,uuid" json:"delegate_id"`
	StartDate  string `validate:"required" json:"start_date"`
	EndDate    string `validate:"required" json:"end_date"`
}
