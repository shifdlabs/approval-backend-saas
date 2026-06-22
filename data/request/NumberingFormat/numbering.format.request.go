package numberingformat

type NumberingFormatRequest struct {
	GroupID          string `validate:"required,uuid" json:"group_id"`
	Name             string `validate:"required,min=1,max=200" json:"name"`
	Format           string `validate:"required" json:"format"`
	Separator        string `validate:"required" json:"separator"`
	IncrementByGroup *bool  `json:"increment_by_group"`
}
