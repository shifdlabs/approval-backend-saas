package document

type RecentDocumentResponse struct {
	ID        string `json:"id"`
	Number    string `json:"number"`
	Subject   string `json:"subject"`
	FromTo    string `json:"from_to"`
	Status    string `json:"status"`
	Type      int    `json:"type"`
	UpdatedAt string `json:"updated_at"`
}
