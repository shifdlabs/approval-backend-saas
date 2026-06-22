package documentnumbers

type DocumentNumbersRequest struct {
	NumberingFormatID string `validate:"required,uuid" json:"numbering_format_id"`
}
