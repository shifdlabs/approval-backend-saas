package aidraft

// DraftVariant is one generated style option. Style is a stable key
// ("formal" | "ringkas" | "semi-formal") — the frontend owns the display
// label/translation for it, this backend does not.
type DraftVariant struct {
	Style      string `json:"style"`
	Draft      string `json:"draft"`
	TokensUsed int    `json:"tokens_used"`
}

// GenerateDraftVariantsResponse carries every style that generated
// successfully. It can hold fewer than len(draftStyles) entries: a single
// style failing does not fail the whole request (see Generate), so the user
// still gets to choose among whatever came back.
type GenerateDraftVariantsResponse struct {
	Variants []DraftVariant `json:"variants"`
}
