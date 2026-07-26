package aidraft

// GenerateDraftRequest is the letter description the user fills in before
// asking the AI for a draft. Limits mirror SAIS's own validation so an
// oversized request is rejected here instead of costing a round trip.
type GenerateDraftRequest struct {
	LetterType string   `json:"letter_type" validate:"required,max=100"`
	Purpose    string   `json:"purpose" validate:"required,max=1000"`
	KeyPoints  []string `json:"key_points" validate:"required,min=1,max=20,dive,required,max=500"`
	Recipient  string   `json:"recipient" validate:"max=255"`
	Tone       string   `json:"tone" validate:"omitempty,oneof=formal semi-formal"`
}
