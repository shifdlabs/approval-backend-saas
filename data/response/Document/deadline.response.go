package document

type DeadlineItemResponse struct {
	ID            string `json:"id"`
	Subject       string `json:"subject"`
	DaysRemaining int    `json:"days_remaining"` // negatif = lewat, 0 = hari ini, positif = sisa hari
}
