package document

type CreateDocumentRequest struct {
	AuthorID              string                      `validate:"required,uuid" json:"authorID"`
	PublicationNumberType int                         `validate:"required,oneof=1 2 3 4" json:"publicationNumberType"` // 1: Auto-Generated, 2: Booking Number, 3: Custom, 4: N/A (No Number)
	PublicationValue      *string                     `json:"publicationValue"`                                        // it could be Booked Number, Format ID or Custom Number
	Type                  int                         `validate:"required,oneof=1 2" json:"type"`                      // 1: Internal, 2: External
	Priority              int                         `validate:"required,oneof=1 2 3" json:"priority"`                // 1: High, 2: Medium, 3: Low
	Subject               string                      `validate:"required,min=1,max=500" json:"subject"`
	Body                  string                      `validate:"required" json:"body"`
	ExternalRecipient     string                      `json:"externalRecipient"`
	Step                  int                         `validate:"required" json:"step"`
	LetterHead            bool                        `json:"letterHead"`
	Status                int                         `json:"status"`
	Recipients            []string                    `json:"recipients"`
	CarbonCopies          []string                    `json:"carbonCopies"`
	Sequences             []DocumentSequence          `json:"sequences"`
	Attachments           []DocumentAttachmentRequest `json:"attachments"`
	References            []string                    `json:"references"`
	TemplateID            *string                     `json:"templateId"`
}

type DocumentSequence struct {
	UserID    string `validate:"required,uuid" json:"userID"`
	Signature bool   `json:"signature"`
}

type DocumentAttachmentRequest struct {
	OriginalName string `validate:"required" json:"originalName"`
	FileName     string `validate:"required" json:"fileName"`
	Path         string `validate:"required" json:"path"`
	Size         string `validate:"required" json:"size"`
	Type         string `validate:"required" json:"type"`
}
