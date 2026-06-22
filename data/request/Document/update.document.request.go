package document

type UpdateDocumentRequest struct {
	Id                    string                      `validate:"required,uuid" json:"id"`
	AuthorID              string                      `validate:"required,uuid" json:"authorID"`
	PublicationNumberType int                         `validate:"required,oneof=1 2 3 4" json:"publicationNumberType"` // 1: Auto-Generated, 2: Booking Number, 3: Custom, 4: N/A (No Number)
	PublicationValue      *string                     `json:"publicationValue"`
	Type                  int                         `validate:"required,oneof=1 2" json:"type"`     // 1: Internal, 2: External
	Priority              int                         `validate:"required,oneof=1 2 3" json:"priority"` // 1: High, 2: Medium, 3: Low
	Subject               string                      `validate:"required,min=1,max=500" json:"subject"`
	Body                  string                      `validate:"required" json:"body"`
	ExternalRecipient     string                      `json:"externalRecipient"`
	LetterHead            bool                        `json:"letterHead"`
	Recipients            []string                    `json:"recipients"`
	CarbonCopies          []string                    `json:"carbonCopies"`
	Sequences             []DocumentSequence          `validate:"required" json:"sequences"`
	NewAttachments        []DocumentAttachmentRequest `json:"newAttachments"`
	IsDraft               bool                        `json:"isDraft"`
	References            []string                    `json:"references"`
}
