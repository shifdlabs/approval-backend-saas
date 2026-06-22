package lettertemplate

type CreateLetterTemplateRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"max=255"`
	Body        string `json:"body"`
}

type UpdateLetterTemplateRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"max=255"`
	Body        string `json:"body"`
}
