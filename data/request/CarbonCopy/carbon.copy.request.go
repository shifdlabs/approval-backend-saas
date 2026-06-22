package position

type CarbonCopyRequest struct {
	DocumentId string   `validate:"required,uuid" json:"documentId"`
	UserIds    []string `validate:"required,min=1,max=100,dive,required" json:"userIds"`
}
