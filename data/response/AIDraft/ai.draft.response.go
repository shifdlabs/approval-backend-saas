package aidraft

// GenerateDraftResponse carries the generated letter body back to the editor.
// TokensUsed is passed through from SAIS so the UI can surface usage if needed;
// billing itself is accounted for in SAIS, not here.
type GenerateDraftResponse struct {
	Draft      string `json:"draft"`
	TokensUsed int    `json:"tokens_used"`
}
