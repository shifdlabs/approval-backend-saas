package position

type AppSettingRequest struct {
	Properties []AppProperty `validate:"required" json:"properties"`
}

type AppProperty struct {
	Key   string `validate:"required" json:"key"`
	Value string `validate:"required" json:"value"`
}
