package document

type ActivityResponse struct {
	ID           string `json:"id"`
	Subject      string `json:"subject"`
	IsApproved   bool   `json:"is_approved"`
	ApproverName string `json:"approver_name"`
	UpdatedAt    string `json:"updated_at"`
}
